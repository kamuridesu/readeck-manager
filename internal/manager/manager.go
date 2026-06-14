package manager

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/kamuridesu/readeck-manager/internal/config"
	"github.com/kamuridesu/readeck-manager/internal/readeck"
	"github.com/kamuridesu/readeck-manager/internal/tagger"
	"github.com/kamuridesu/readeck-manager/internal/utils"
)

type Manager struct {
	cfg      *config.Config
	rd       *readeck.Readeck
	aiTagger *tagger.Tagger
}

func New(cfg *config.Config) *Manager {
	return &Manager{
		cfg:      cfg,
		rd:       readeck.New(cfg),
		aiTagger: tagger.New(cfg),
	}
}

func (m *Manager) TagUntaggedLinks(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	slog.Info("Loading untagged links...")
	untagged, err := m.rd.GetBookmarks(ctx, readeck.BookmarkListOptions{
		HasLabels: new(false),
	})
	if err != nil {
		return err
	}
	if len(untagged) == 0 {
		slog.Info("No untagged bookmarks found")
		return nil
	}
	slog.Info("Loading labels...")
	labels, err := m.rd.GetLabels(ctx)
	strLabels := []string{}
	for _, label := range labels {
		strLabels = append(strLabels, label.Name)
	}
	if err != nil {
		return err
	}
	slog.Info("Starting tagging process...")

	var wg sync.WaitGroup
	const maxConcurrentTasks = 5
	sem := make(chan struct{}, maxConcurrentTasks)
	for _, b := range untagged {
		bookmark := b

		sem <- struct{}{}

		wg.Go(func() {
			defer func() {
				<-sem
			}()

			title := bookmark.Title
			if len(title) > 30 {
				title = title[:30]
			} else if len(title) < 30 {
				title = title + strings.Repeat(" ", 30-len(title))
			}

			slog.Info(fmt.Sprintf("[%s] Fetching contents for bookmark...", title))
			content, err := m.rd.GetBookmarkHTML(ctx, bookmark.ID)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			content, err = utils.StripHTML(content)
			if err != nil {
				return
			}
			if len(content) > 2000 {
				content = content[:2000]
			}
			slog.Info(fmt.Sprintf("[%s] Generating labels...", title))
			aiTags, err := m.aiTagger.GenerateLabels(ctx, &bookmark, strLabels, content)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			slog.Info(fmt.Sprintf("[%s] Tagging bookmark with labels '%s'...", title, strings.Join(aiTags, ", ")))
			err = m.rd.UpdateBookmarkLabels(ctx, bookmark.ID, aiTags)
			if err != nil {
				slog.Info(err.Error())
				return
			}
			slog.Info(fmt.Sprintf("[%s] Bookmark tagged!", title))

		})
	}

	wg.Wait()

	return nil
}

func (m *Manager) DeletedTaggedWithDelete(ctx context.Context, dryRun bool) error {
	toDelete, err := m.rd.GetBookmarks(ctx, readeck.BookmarkListOptions{
		Labels: []string{"DELETE"},
	})
	if err != nil {
		return err
	}

	for _, bookmark := range toDelete {
		fmt.Printf("Deleting '%s'\n", bookmark.Title)
		if dryRun {
			continue
		}
		err := m.rd.DeleteBookmark(ctx, bookmark.ID)
		if err != nil {
			slog.Error(err.Error())
			continue
		}
	}

	return nil
}

func (m *Manager) GetBrokenBookmarks(ctx context.Context) error {
	bms, err := m.rd.GetBookmarks(ctx, readeck.BookmarkListOptions{
		HasErrors: new(true),
	})
	if err != nil {
		return err
	}

	maxBMTitleChars := 30

	for i, bm := range bms {
		title := bm.Title
		if len(title) > maxBMTitleChars {
			title = title[:maxBMTitleChars]
		} else {
			title = title + strings.Repeat(" ", maxBMTitleChars-len(title))
		}

		url := strings.Replace(bm.Href, "/api", "", 1)

		fmt.Printf("%d. '%s' has errors: %s\n", i+1, title, url)
	}
	return nil
}

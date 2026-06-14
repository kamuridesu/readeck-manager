package readeck

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kamuridesu/readeck-manager/internal/config"
	"github.com/kamuridesu/readeck-manager/internal/request"
)

func New(cfg *config.Config) *Readeck {
	return &Readeck{
		baseUrl: cfg.ReadeckURL,
		token:   cfg.ReadeckToken,
		client:  http.DefaultClient,
	}
}

func (r *Readeck) getUrl(toReplace string) string {
	return fmt.Sprintf(toReplace, r.baseUrl)
}

func (r *Readeck) GetBookmarks(ctx context.Context, options BookmarkListOptions) ([]Bookmark, error) {
	var bookmarks []Bookmark

	u, err := url.Parse(r.getUrl("%s/bookmarks"))
	if err != nil {
		return nil, err
	}

	q := u.Query()

	if options.Limit > 0 {
		q.Add("limit", strconv.Itoa(options.Limit))
	}
	if options.Offset > 0 {
		q.Add("offset", strconv.Itoa(options.Offset))
	}
	for _, sort := range options.Sort {
		q.Add("sort", string(sort))
	}
	if options.Search != "" {
		q.Add("search", options.Search)
	}
	if options.Title != "" {
		q.Add("title", options.Title)
	}
	if options.Author != "" {
		q.Add("author", options.Author)
	}
	if options.Site != "" {
		q.Add("site", options.Site)
	}
	for _, t := range options.Type {
		q.Add("type", string(t))
	}
	for _, label := range options.Labels {
		q.Add("labels", label)
	}

	if options.IsLoaded != nil {
		q.Add("is_loaded", strconv.FormatBool(*options.IsLoaded))
	}
	if options.HasErrors != nil {
		q.Add("has_errors", strconv.FormatBool(*options.HasErrors))
	}
	if options.HasLabels != nil {
		q.Add("has_labels", strconv.FormatBool(*options.HasLabels))
	}
	if options.IsMarked != nil {
		q.Add("is_marked", strconv.FormatBool(*options.IsMarked))
	}
	if options.IsArchived != nil {
		q.Add("is_archived", strconv.FormatBool(*options.IsArchived))
	}

	if options.RangeStart != "" {
		q.Add("range_start", options.RangeStart)
	}
	if options.RangeEnd != "" {
		q.Add("range_end", options.RangeEnd)
	}
	for _, rs := range options.ReadStatus {
		q.Add("read_status", string(rs))
	}
	if options.Id != "" {
		q.Add("id", options.Id)
	}
	if options.Collection != "" {
		q.Add("collection", options.Collection)
	}

	u.RawQuery = q.Encode()

	headers := map[string]string{
		"Authorization": "Bearer " + r.token,
		"Accept":        "application/json",
	}

	err = request.SendGETRequest(ctx, r.client, u.String(), &bookmarks, headers)
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func (r *Readeck) GetLabels(ctx context.Context) ([]Label, error) {
	var labels []Label
	headers := map[string]string{"Authorization": "Bearer " + r.token, "Accept": "application/json"}
	err := request.SendGETRequest(ctx, r.client, r.getUrl("%s/bookmarks/labels"), &labels, headers)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (r *Readeck) GetBookmarkHTML(ctx context.Context, id string) (string, error) {
	var html string
	headers := map[string]string{"Authorization": "Bearer " + r.token, "Accept": "text/html"}
	err := request.SendGETRequest(ctx, r.client, fmt.Sprintf("%s/bookmarks/%s/article", r.baseUrl, id), &html, headers)
	if err != nil {
		return "", err
	}
	return html, nil
}

func (r *Readeck) UpdateBookmarkLabels(ctx context.Context, id string, labels []string) error {
	payload := map[string][]string{"add_labels": labels}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to parse payload to json: %e", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", fmt.Sprintf("%s/bookmarks/%s", r.baseUrl, id), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to build request: %e", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := r.client.Do(req)
	if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("an unexpected error happened")
	}
	defer res.Body.Close()
	return nil
}

func (r *Readeck) DeleteBookmark(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/bookmarks/%s", r.baseUrl, id), nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %e", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := r.client.Do(req)
	if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("an unexpected error happened")
	}
	defer res.Body.Close()
	return nil
}

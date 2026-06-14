package readeck

import (
	"net/http"
	"time"
)

type Readeck struct {
	baseUrl string
	token   string
	client  *http.Client
}

type Bookmark struct {
	ID            string     `json:"id"`
	Href          string     `json:"href"`
	Created       time.Time  `json:"created"`
	Updated       time.Time  `json:"updated"`
	State         int        `json:"state"`
	Loaded        bool       `json:"loaded"`
	URL           string     `json:"url"`
	Title         string     `json:"title"`
	SiteName      string     `json:"site_name"`
	Site          string     `json:"site"`
	Published     *time.Time `json:"published"`
	Authors       []string   `json:"authors"`
	Lang          string     `json:"lang"`
	TextDirection string     `json:"text_direction"`
	DocumentType  string     `json:"document_type"`
	Type          string     `json:"type"`
	HasArticle    bool       `json:"has_article"`
	Description   string     `json:"description"`
	IsDeleted     bool       `json:"is_deleted"`
	IsMarked      bool       `json:"is_marked"`
	IsArchived    bool       `json:"is_archived"`
	ReadProgress  int        `json:"read_progress"`
	Labels        []string   `json:"labels"`
	WordCount     int        `json:"word_count"`
	ReadingTime   int        `json:"reading_time"`
	Resources     Resources  `json:"resources"`
}

type Resources struct {
	Log   ResourceSrc `json:"log"`
	Props ResourceSrc `json:"props"`

	Article   *ResourceSrc   `json:"article,omitempty"`
	Icon      *ImageResource `json:"icon,omitempty"`
	Image     *ImageResource `json:"image,omitempty"`
	Thumbnail *ImageResource `json:"thumbnail,omitempty"`
}

type ResourceSrc struct {
	Src string `json:"src"`
}

type ImageResource struct {
	Src    string `json:"src"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Label struct {
	Name          string `json:"name"`
	Count         int    `json:"count"`
	HRef          string `json:"href"`
	HRefBookmarks string `json:"href_bookmarks"`
}

type BookmarkListOptions struct {
	Limit      int
	Offset     int
	Sort       []Sort
	Search     string
	Title      string
	Author     string
	Site       string
	Type       []Type
	Labels     []string
	IsLoaded   *bool
	HasErrors  *bool
	HasLabels  *bool
	IsMarked   *bool
	IsArchived *bool
	RangeStart string
	RangeEnd   string
	ReadStatus []ReadStatus
	Id         string
	Collection string
}

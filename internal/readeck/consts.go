package readeck

type Sort = string

const (
	CreatedAsc    Sort = "created"
	CreatedDesc   Sort = "-created"
	DomainAsc     Sort = "domain"
	DomainDesc    Sort = "-domain"
	DurationAsc   Sort = "duration"
	DurationDesc  Sort = "-duration"
	PublishedAsc  Sort = "published"
	PublishedDesc Sort = "-published"
	SiteAsc       Sort = "site"
	SiteDesc      Sort = "-site"
	TitleAsc      Sort = "title"
	TitleDesc     Sort = "-title"
)

type Type = string

const (
	Article Type = "article"
	Photo   Type = "photo"
	Video   Type = "video"
)

type ReadStatus = string

const (
	Unread  ReadStatus = "unread"
	Reading ReadStatus = "reading"
	Read    ReadStatus = "read"
)

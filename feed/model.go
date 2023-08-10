package feed

import "time"

type NewsletterNewsItem struct {
	ArticleURL        string
	NewsArticleID     int
	PublishDate       time.Time
	Taxonomies        string
	TeaserText        string
	Subtitle          string
	ThumbnailImageURL string
	Title             string
	BodyText          string
	GalleryImageURLs  string
	VideoURL          string
	OptaMatchId       string
	LastUpdateDate    time.Time
	IsPublished       bool
}

type newListInformation struct {
	NewsletterNewsItems []NewsletterNewsItem `xml:"NewsletterNewsItems>NewsletterNewsItem"`
}

type fetchManyResp struct {
	NewListInformation newListInformation `xml:"NewListInformation"`
}

type fetchOneResp struct {
	NewsArticle NewsletterNewsItem `xml:"NewsArticle"`
}

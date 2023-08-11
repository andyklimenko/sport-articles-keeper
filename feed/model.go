package feed

import (
	"encoding/xml"
)

type NewsletterNewsItem struct {
	ArticleURL        string `xml:"ArticleURL"`
	NewsArticleID     int    `xml:"NewsArticleID"`
	PublishDate       string `xml:"PublishDate"`
	Taxonomies        string `xml:"Taxonomies"`
	TeaserText        string `xml:"TeaserText"`
	Subtitle          string `xml:"Subtitle"`
	ThumbnailImageURL string `xml:"ThumbnailImageURL"`
	Title             string `xml:"Title"`
	BodyText          string `xml:"BodyText"`
	GalleryImageURLs  string `xml:"GalleryImageURLs"`
	VideoURL          string `xml:"VideoURL"`
	OptaMatchId       string `xml:"OptaMatchId"`
	LastUpdateDate    string `xml:"LastUpdateDate"`
	IsPublished       bool   `xml:"IsPublished"`
}

type fetchManyResp struct {
	XMLName             xml.Name             `xml:"NewListInformation"`
	NewsletterNewsItems []NewsletterNewsItem `xml:"NewsletterNewsItems>NewsletterNewsItem"`
}

type fetchOneResp struct {
	XMLName     xml.Name           `xml:"NewsArticleInformation"`
	NewsArticle NewsletterNewsItem `xml:"NewsArticle"`
}

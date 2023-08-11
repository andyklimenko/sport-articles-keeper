package consumer

import (
	"github.com/andyklimenko/sport-articles-keeper/feed"
	"github.com/andyklimenko/sport-articles-keeper/model"
)

func articleFromFeed(item feed.NewsletterNewsItem) model.Article {
	return model.Article{
		ID:            "",
		NewsArticleID: item.NewsArticleID,
		TeamID:        "",
		OptaMatchId:   item.OptaMatchId,
		Title:         item.Title,
		Type:          nil,
		Teaser:        item.TeaserText,
		Content:       item.BodyText,
		Url:           item.ArticleURL,
		ImageUrl:      item.ThumbnailImageURL,
		GalleryUrls:   item.GalleryImageURLs,
		VideoUrl:      item.VideoURL,
		Published:     item.PublishDate,
	}
}

package storage

import (
	"github.com/andyklimenko/sport-articles-keeper/model"
)

type article struct {
	ID            string      `bson:"id"`
	NewsArticleID int         `bson:"newsArticleId"`
	TeamID        string      `bson:"teamId"`
	OptaMatchId   interface{} `bson:"optaMatchId"`
	Title         string      `bson:"title"`
	Subtitle      string      `bson:"subtitle"`
	Type          []string    `bson:"type"`
	Teaser        string      `bson:"teaser"`
	Content       string      `bson:"content"`
	Url           string      `bson:"url"`
	ImageUrl      string      `bson:"imageUrl"`
	GalleryUrls   string      `bson:"galleryUrls"`
	VideoUrl      string      `bson:"videoUrl"`
	Published     string      `bson:"published"`
}

func (a article) model() model.Article {
	return model.Article{
		ID:            a.ID,
		NewsArticleID: a.NewsArticleID,
		TeamID:        a.TeamID,
		OptaMatchId:   a.OptaMatchId,
		Title:         a.Title,
		Subtitle:      a.Subtitle,
		Type:          a.Type,
		Teaser:        a.Teaser,
		Content:       a.Content,
		Url:           a.Url,
		ImageUrl:      a.ImageUrl,
		GalleryUrls:   a.GalleryUrls,
		VideoUrl:      a.VideoUrl,
		Published:     a.Published,
	}
}

func articleFromModel(a model.Article) article {
	return article{
		ID:            a.ID,
		NewsArticleID: a.NewsArticleID,
		TeamID:        a.TeamID,
		OptaMatchId:   a.OptaMatchId,
		Title:         a.Title,
		Type:          a.Type,
		Teaser:        a.Teaser,
		Content:       a.Content,
		Url:           a.Url,
		ImageUrl:      a.ImageUrl,
		GalleryUrls:   a.GalleryUrls,
		VideoUrl:      a.VideoUrl,
		Published:     a.Published,
	}
}

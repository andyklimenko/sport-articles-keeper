package model

import "time"

type Article struct {
	ID            string      `json:"id"`
	NewsArticleID int         `json:"-"`
	TeamID        string      `json:"teamId"`
	OptaMatchId   interface{} `json:"optaMatchId"`
	Title         string      `json:"title"`
	Type          []string    `json:"type"`
	Teaser        string      `json:"teaser"`
	Content       string      `json:"content"`
	Url           string      `json:"url"`
	ImageUrl      string      `json:"imageUrl"`
	GalleryUrls   interface{} `json:"galleryUrls"`
	VideoUrl      interface{} `json:"videoUrl"`
	Published     time.Time   `json:"published"`
}

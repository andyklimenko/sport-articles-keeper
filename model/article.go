package model

type Article struct {
	ID            string      `json:"id"`
	NewsArticleID int         `json:"-"`
	TeamID        string      `json:"teamId"`
	OptaMatchId   interface{} `json:"optaMatchId"`
	Title         string      `json:"title"`
	Subtitle      string      `json:"subtitle"`
	Type          []string    `json:"type"`
	Teaser        string      `json:"teaser"`
	Content       string      `json:"content"`
	Url           string      `json:"url"`
	ImageUrl      string      `json:"imageUrl"`
	GalleryUrls   string      `json:"galleryUrls"`
	VideoUrl      string      `json:"videoUrl"`
	Published     string      `json:"published"`
}

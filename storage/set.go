package storage

import (
	"context"

	"github.com/andyklimenko/sport-articles-keeper/model"
)

func (s *Storage) InsertMany(ctx context.Context, documents []model.Article) error {
	in := make([]interface{}, 0, len(documents))
	for _, d := range documents {
		in = append(in, articleFromModel(d))
	}

	_, err := s.newsletterCollection.InsertMany(ctx, in)
	return err
}

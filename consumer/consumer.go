package consumer

import (
	"context"
	"log/slog"

	"github.com/andyklimenko/sport-articles-keeper/feed"
	"github.com/andyklimenko/sport-articles-keeper/model"
	"github.com/andyklimenko/sport-articles-keeper/storage"
)

type repo interface {
	GetMany(ctx context.Context, filter storage.ArticleFilter) ([]model.Article, error)
	InsertMany(ctx context.Context, articles []model.Article) error
}

type fetcher interface {
	FetchOne(ctx context.Context, id int) (feed.NewsletterNewsItem, error)
}

type Consumer struct {
	repo    repo
	fetcher fetcher
}

func (c *Consumer) consume(consumeCh <-chan []feed.NewsletterNewsItem) {
	for batch := range consumeCh {
		slog.With(slog.Any("size", len(batch))).Info("got new batch")

		if err := c.onBatch(batch); err != nil {
			slog.Error("processing batch", slog.Any("error", err))
		}
	}
}

func Start(r repo, f fetcher, consumeCh <-chan []feed.NewsletterNewsItem) {
	c := Consumer{
		repo:    r,
		fetcher: f,
	}

	go c.consume(consumeCh)
}

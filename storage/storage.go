package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/andyklimenko/sport-articles-keeper/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const newsCollectionName = "newsletterNews"

type Storage struct {
	newsletterCollection *mongo.Collection

	cli *mongo.Client
}

func (s *Storage) Disconnect(ctx context.Context) error {
	return s.cli.Disconnect(ctx)
}

func New(cfg config.Storage) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	db := client.Database(cfg.DbName)

	return &Storage{
		cli:                  client,
		newsletterCollection: db.Collection(newsCollectionName),
	}, nil
}

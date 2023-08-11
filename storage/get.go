package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/andyklimenko/sport-articles-keeper/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleFilter struct {
	ArticleIDs []int // to be used internally
	Limit      int64
	Offset     int64
}

func (f ArticleFilter) opts() *options.FindOptions {
	res := options.Find()
	if f.Offset > 0 {
		res = res.SetSkip(f.Offset)
	}

	if f.Limit > 0 {
		res = res.SetLimit(f.Limit)
	}

	return res
}

func (f ArticleFilter) m() bson.M {
	res := bson.M{}
	if len(f.ArticleIDs) > 0 {
		res["newsArticleId"] = map[string][]int{
			"$in": f.ArticleIDs,
		}
	}

	return res
}

func (s *Storage) GetMany(ctx context.Context, filter ArticleFilter) ([]model.Article, int64, error) {
	f := filter.m()

	cur, err := s.newsletterCollection.Find(ctx, f, filter.opts())
	if err != nil {
		return nil, 0, fmt.Errorf("looking in db: %w", err)
	}

	var out []article
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, fmt.Errorf("decode result: %w", err)
	}

	if len(out) == 0 {
		return nil, 0, ErrNotFound
	}

	total, err := s.newsletterCollection.CountDocuments(ctx, f, options.Count().SetHint("_id_"))
	if err != nil {
		return nil, 0, fmt.Errorf("counting total documents number: %w", err)
	}

	res := make([]model.Article, 0, len(out))
	for _, a := range out {
		res = append(res, a.model())
	}

	return res, total, nil
}

func (s *Storage) GetOne(ctx context.Context, id string) (model.Article, error) {
	res := s.newsletterCollection.FindOne(ctx, bson.M{"id": id})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Article{}, ErrNotFound
	}

	if err != nil {
		return model.Article{}, fmt.Errorf("looking for article: %w", err)
	}

	var a article
	if err := res.Decode(&a); err != nil {
		return model.Article{}, fmt.Errorf("looking for article by ID=%s: %w", id, err)
	}

	return a.model(), nil
}

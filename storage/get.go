package storage

import (
	"context"
	"fmt"

	"github.com/andyklimenko/sport-articles-keeper/model"
	"go.mongodb.org/mongo-driver/bson"
)

type ArticleFilter struct {
	ArticleIDs []int // to be used internally
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

func (s *Storage) GetMany(ctx context.Context, filter ArticleFilter) ([]model.Article, error) {
	cur, err := s.newsletterCollection.Find(ctx, filter.m())
	if err != nil {
		return nil, fmt.Errorf("looking in db: %w", err)
	}

	var out []article
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode result: %w", err)
	}

	if len(out) == 0 {
		return nil, ErrNotFound
	}

	res := make([]model.Article, 0, len(out))
	for _, a := range out {
		res = append(res, a.model())
	}

	return res, nil
}

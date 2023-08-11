package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/andyklimenko/sport-articles-keeper/feed"
	"github.com/andyklimenko/sport-articles-keeper/model"
	"github.com/andyklimenko/sport-articles-keeper/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mckRepo struct {
	mock.Mock
	insertManyCh chan []model.Article
}

func (r *mckRepo) GetMany(ctx context.Context, filter storage.ArticleFilter) ([]model.Article, int64, error) {
	res := r.Called(ctx, filter)
	return res.Get(0).([]model.Article), int64(res.Int(1)), res.Error(2)
}

func (r *mckRepo) InsertMany(ctx context.Context, articles []model.Article) error {
	err := r.Called(ctx, articles).Error(0)
	go func() {
		r.insertManyCh <- articles
	}()

	return err
}

type mckFetcher struct {
	mock.Mock
}

func (m *mckFetcher) FetchOne(ctx context.Context, id int) (feed.NewsletterNewsItem, error) {
	res := m.Called(ctx, id)
	return res.Get(0).(feed.NewsletterNewsItem), res.Error(1)
}

func TestConsumer(t *testing.T) {
	t.Parallel()

	r := mckRepo{
		insertManyCh: make(chan []model.Article),
	}
	f := mckFetcher{}
	consumeCh := make(chan []feed.NewsletterNewsItem)

	Start(&r, &f, consumeCh)

	articlesWeAlreadyHave := []model.Article{
		{
			ID:            "1",
			NewsArticleID: 1,
			Title:         "abc",
		},
		{
			ID:            "2",
			NewsArticleID: 2,
			Title:         "def",
		},
	}

	newArticle := feed.NewsletterNewsItem{
		NewsArticleID: 3,
		Title:         "xyz",
	}

	articlesWeFetched := []feed.NewsletterNewsItem{
		{
			NewsArticleID: 1,
			Title:         "abc",
		},
		{
			NewsArticleID: 2,
			Title:         "abc",
		},
		newArticle,
	}

	expectedFilter := storage.ArticleFilter{
		ArticleIDs: []int{1, 2, 3},
	}

	r.On("GetMany", mock.Anything, expectedFilter).Return(articlesWeAlreadyHave, 2, nil).Once()
	r.On("InsertMany",
		mock.Anything,
		mock.MatchedBy(func(in []model.Article) bool {
			if len(in) != 1 {
				return false
			}

			return in[0].NewsArticleID == 3 && in[0].Title == "xyz"
		}),
	).Return(nil).Once()
	f.On("FetchOne", mock.Anything, 3).Return(newArticle, nil).Once()

	consumeCh <- articlesWeFetched

	select {
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	case savedArticles := <-r.insertManyCh:
		require.Len(t, savedArticles, 1)
		assert.Equal(t, 3, savedArticles[0].NewsArticleID)
		assert.Equal(t, "xyz", savedArticles[0].Title)
	}

	close(consumeCh)

	r.AssertExpectations(t)
	f.AssertExpectations(t)
}

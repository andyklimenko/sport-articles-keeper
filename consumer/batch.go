package consumer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/andyklimenko/sport-articles-keeper/feed"
	"github.com/andyklimenko/sport-articles-keeper/model"
	"github.com/andyklimenko/sport-articles-keeper/storage"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func (c *Consumer) onBatch(batch []feed.NewsletterNewsItem) error {
	if len(batch) == 0 {
		return nil
	}

	articleIDs := make([]int, 0, len(batch))
	incomingArticlesMap := make(map[int]feed.NewsletterNewsItem, len(batch))
	for _, b := range batch {
		articleIDs = append(articleIDs, b.NewsArticleID)
		incomingArticlesMap[b.NewsArticleID] = b
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := storage.ArticleFilter{
		ArticleIDs: articleIDs,
	}
	existingArticles, _, err := c.repo.GetMany(ctx, filter)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("looking for existing articles: %w", err)
	}

	existingArticleMap := make(map[int]model.Article, len(existingArticles))
	for _, item := range existingArticles {
		existingArticleMap[item.NewsArticleID] = item
	}

	var articleIDsToGetMoreDetails []int
	newArticlesMap := map[int]model.Article{}
	for articleID, feedItem := range incomingArticlesMap {
		if _, found := existingArticleMap[articleID]; found {
			continue // the document was already saved, skipping it
		}

		articleIDsToGetMoreDetails = append(articleIDsToGetMoreDetails, feedItem.NewsArticleID)
		newArticlesMap[feedItem.NewsArticleID] = articleFromFeed(feedItem)
	}

	newsletterNewsItemsToSave := make([]model.Article, 0, len(newArticlesMap))
	articlesCh := c.getMoreDetailsAsync(articleIDsToGetMoreDetails)
	for article := range articlesCh {
		a, found := newArticlesMap[article.NewsArticleID]
		if !found {
			slog.With("articleID", article.NewsArticleID).Warn("unexpected article occurred!!")
			continue
		}

		a.ID = uuid.NewString()
		a.Content = article.BodyText
		a.GalleryUrls = article.GalleryImageURLs
		a.VideoUrl = article.VideoURL
		newsletterNewsItemsToSave = append(newsletterNewsItemsToSave, a)
	}

	if err := c.saveNewArticles(ctx, newsletterNewsItemsToSave); err != nil {
		return fmt.Errorf("saving new articles: %w", err)
	}

	slog.With("article batch size", len(newsletterNewsItemsToSave)).Info("processed")

	return nil
}

func (c *Consumer) saveNewArticles(ctx context.Context, articles []model.Article) error {
	if len(articles) == 0 {
		return nil
	}

	return c.repo.InsertMany(ctx, articles)
}

func (c *Consumer) getMoreDetailsAsync(ids []int) <-chan feed.NewsletterNewsItem {
	outCh := make(chan feed.NewsletterNewsItem)

	go func() {
		defer close(outCh)

		g, _ := errgroup.WithContext(context.Background())
		for i := range ids {
			id := ids[i]

			g.Go(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				got, err := c.fetcher.FetchOne(ctx, id)
				if err != nil {
					return fmt.Errorf("get single newsletter item detail: %w", err)
				}

				outCh <- got
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			slog.With("article IDs", ids).Error("get articles detailed info", slog.Any("error", err))
		}
	}()
	return outCh
}

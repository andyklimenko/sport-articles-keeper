package poller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/andyklimenko/sport-articles-keeper/feed"
	"github.com/go-co-op/gocron"
)

type fetcher interface {
	FetchMany(ctx context.Context) ([]feed.NewsletterNewsItem, error)
}

type Poller struct {
	fetcher           fetcher
	scheduler         *gocron.Scheduler
	newsletterBatchCh chan<- []feed.NewsletterNewsItem
}

func (p *Poller) StartAsync() {
	p.scheduler.StartAsync()
}

func (p *Poller) Stop() {
	p.scheduler.Stop()
	close(p.newsletterBatchCh)
}

func (p *Poller) poll() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	batchGot, err := p.fetcher.FetchMany(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to query the feed", slog.Any("error", err))
	}

	go func() {
		p.newsletterBatchCh <- batchGot
	}()
}

func New(interval time.Duration, f fetcher, newsletterBatchCh chan<- []feed.NewsletterNewsItem) (*Poller, error) {
	p := Poller{
		fetcher:           f,
		newsletterBatchCh: newsletterBatchCh,
		scheduler:         gocron.NewScheduler(time.UTC),
	}

	_, err := p.scheduler.Every(interval).Do(p.poll)
	if err != nil {
		return nil, fmt.Errorf("initiate periodic job execution: %w", err)
	}

	return &p, nil
}

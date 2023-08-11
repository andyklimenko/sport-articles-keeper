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
	newsletterBatchCh chan []feed.NewsletterNewsItem
}

// StartBlocking should be call to start polling
func (p *Poller) StartBlocking() {
	p.scheduler.StartBlocking()
}

// Stop must be called for graceful shutdown sake
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

// New returns the poller and the channel to be used by consumer for reading the data fetched. The channel will be closed automatically during graceful shutdown.
func New(interval time.Duration, f fetcher) (*Poller, <-chan []feed.NewsletterNewsItem, error) {
	p := Poller{
		fetcher:           f,
		newsletterBatchCh: make(chan []feed.NewsletterNewsItem),
		scheduler:         gocron.NewScheduler(time.UTC),
	}

	_, err := p.scheduler.Every(interval).Do(p.poll)
	if err != nil {
		return nil, nil, fmt.Errorf("initiate periodic job execution: %w", err)
	}

	return &p, p.newsletterBatchCh, nil
}

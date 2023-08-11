package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andyklimenko/sport-articles-keeper/config"
	"github.com/andyklimenko/sport-articles-keeper/consumer"
	"github.com/andyklimenko/sport-articles-keeper/feed"
	"github.com/andyklimenko/sport-articles-keeper/poller"
	"github.com/andyklimenko/sport-articles-keeper/storage"
)

func main() {
	var cfg config.Config
	if err := cfg.Load(); err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	repo, err := storage.New(cfg.Storage)
	if err != nil {
		panic(fmt.Errorf("init storage: %w", err))
	}

	feedClient := feed.New(cfg.Feed.ApiURL, cfg.Feed.BatchSize)

	p, consumeCh, err := poller.New(cfg.Poller.Interval, feedClient)
	if err != nil {
		panic(fmt.Errorf("init poller: %w", err))
	}

	consumer.Start(repo, feedClient, consumeCh)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-stopCh
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		p.Stop()

		if err := repo.Disconnect(ctx); err != nil {
			slog.Error("disconnect db", slog.Any("error", err))
		}

		slog.Info("graceful shutdown completed")
	}()

	slog.Info("Starting polling..")
	p.StartBlocking()
}

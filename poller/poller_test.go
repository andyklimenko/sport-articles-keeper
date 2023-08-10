package poller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andyklimenko/sport-articles-keeper/feed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mckFetcher struct {
	mock.Mock
}

func (m *mckFetcher) FetchMany(ctx context.Context) ([]feed.NewsletterNewsItem, error) {
	res := m.Called(ctx)
	return res.Get(0).([]feed.NewsletterNewsItem), res.Error(1)
}

func TestPollingNoErrors(t *testing.T) {
	var f mckFetcher

	const interval = 100 * time.Millisecond
	p, readCh, err := New(interval, &f)
	require.NoError(t, err)

	expected := []feed.NewsletterNewsItem{
		{
			ArticleURL:        "https://www.htafc.com/news/2023/august/p4p11-raises-over-73000-for-charity/",
			NewsArticleID:     612072,
			PublishDate:       time.Date(2023, 8, 9, 15, 0, 0, 0, time.UTC),
			Taxonomies:        "Community",
			TeaserText:        "Yorkshire Air Ambulance, Huddersfield Town Foundation, Andy&#8217;s Man Club and Ruddi&#8217;s Retreat receive over £18,000 each",
			ThumbnailImageURL: "https://www.htafc.com/api/image/feedassets/f81c3def-ba05-4e1b-bd46-7f499c6def88/Medium/p4p-cheque-2023.png",
			Title:             "P4P11 RAISES OVER £73,000 FOR CHARITY!",
			LastUpdateDate:    time.Date(2023, 8, 9, 15, 0, 10, 0, time.UTC),
			IsPublished:       true,
		},
	}

	f.On("FetchMany", mock.Anything).Return(expected, nil)

	p.StartAsync()

	select {
	case <-time.After(2 * interval):
		t.Fatal("timeout")
	case got := <-readCh:
		assert.Equal(t, expected, got)
	}

	p.Stop()

	_, opened := <-readCh
	require.False(t, opened)

	f.AssertExpectations(t)
}

func TestPollingWithErrors(t *testing.T) {
	var f mckFetcher

	const interval = 100 * time.Millisecond
	p, readCh, err := New(interval, &f)
	require.NoError(t, err)

	f.On("FetchMany", mock.Anything).Return([]feed.NewsletterNewsItem{}, errors.New("horrible failure occurred"))

	p.StartAsync()

	select {
	case <-time.After(2 * interval):
		t.Fatal("timeout")
	case got := <-readCh:
		assert.Empty(t, got, got)
	}

	p.Stop()

	_, opened := <-readCh
	require.False(t, opened)

	f.AssertExpectations(t)
}

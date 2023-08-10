package feed

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedNewsletterItems = []NewsletterNewsItem{
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
	{
		ArticleURL:        "https://www.htafc.com/news/2023/august/john-coddington-1937--2023/",
		NewsArticleID:     612074,
		PublishDate:       time.Date(2023, 8, 9, 9, 5, 1, 0, time.UTC),
		Taxonomies:        "Club News",
		TeaserText:        "Long serving Town defender passes away aged 85",
		ThumbnailImageURL: "https://www.htafc.com/api/image/feedassets/7c5d907f-d460-4f1b-bc46-46283c789608/Medium/rip-johncoddington-16x9.jpg",
		Title:             "JOHN CODDINGTON: 1937 &#8211; 2023",
		LastUpdateDate:    time.Date(2023, 8, 9, 9, 8, 0, 0, time.UTC),
		IsPublished:       true,
	},
}

func TestFetchMany(t *testing.T) {
	t.Parallel()

	const expectedBatchSize = 10

	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			qParams := r.URL.Query()
			count := qParams.Get("count")
			if count != strconv.Itoa(expectedBatchSize) {
				w.WriteHeader(http.StatusBadRequest)
				t.Errorf("unexpected 'count' value %s", count)
				return
			}

			w.Header().Set("Content-Type", "application/xml")

			resp := fetchManyResp{
				NewListInformation: newListInformation{
					NewsletterNewsItems: expectedNewsletterItems,
				},
			}
			if err := xml.NewEncoder(w).Encode(resp); err != nil {
				t.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}),
	)

	defer srv.Close()

	parsedURL, err := url.Parse(srv.URL)
	require.NoError(t, err)

	f := New(*parsedURL, 10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newsletterItemsGot, err := f.FetchMany(ctx)
	require.NoError(t, err)
	require.Len(t, newsletterItemsGot, 2)

	assert.Equal(t, expectedNewsletterItems, newsletterItemsGot)
}

func TestFetchOne(t *testing.T) {
	expectedID := expectedNewsletterItems[0].NewsArticleID

	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			qParams := r.URL.Query()
			id := qParams.Get("id")
			if id != strconv.Itoa(expectedID) {
				w.WriteHeader(http.StatusBadRequest)
				t.Errorf("unexpected 'id' value %s", id)
				return
			}

			w.Header().Set("Content-Type", "application/xml")

			resp := fetchOneResp{
				NewsArticle: expectedNewsletterItems[0],
			}
			if err := xml.NewEncoder(w).Encode(resp); err != nil {
				t.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}),
	)
	defer srv.Close()

	parsedURL, err := url.Parse(srv.URL)
	require.NoError(t, err)

	f := New(*parsedURL, -1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	articleGot, err := f.FetchOne(ctx, expectedID)
	require.NoError(t, err)

	assert.Equal(t, expectedNewsletterItems[0], articleGot)
}

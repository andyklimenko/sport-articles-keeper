package feed

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type fetcher struct {
	cli        *http.Client
	apiURL     url.URL
	getManyURL url.URL
}

func (f *fetcher) getRequest(ctx context.Context, reqURL string) (*http.Response, func(), error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, func() {}, fmt.Errorf("init request: %w", err)
	}

	resp, err := f.cli.Do(r)
	if err != nil {
		return nil, func() {}, fmt.Errorf("execute request: %w", err)
	}

	closerFunc := func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("can't close response body: ", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	return resp, closerFunc, nil
}

func (f *fetcher) FetchMany(ctx context.Context) ([]NewsletterNewsItem, error) {
	resp, closer, err := f.getRequest(ctx, f.getManyURL.String())
	defer closer()

	if err != nil {
		return nil, err
	}

	var respBody fetchManyResp
	if err := xml.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return respBody.NewListInformation.NewsletterNewsItems, nil
}

func (f *fetcher) FetchOne(ctx context.Context, id int) (NewsletterNewsItem, error) {
	getByIdURL := f.apiURL
	qParams := getByIdURL.Query()
	qParams.Add("id", strconv.Itoa(id))
	getByIdURL.RawQuery = qParams.Encode()

	resp, closer, err := f.getRequest(ctx, getByIdURL.String())
	defer closer()

	if err != nil {
		return NewsletterNewsItem{}, err
	}

	var respBody fetchOneResp
	if err := xml.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return NewsletterNewsItem{}, fmt.Errorf("decode response body: %w", err)
	}

	return respBody.NewsArticle, nil
}

func newFetcher(apiURL url.URL, batchSize int) *fetcher {
	//we're going to request multiple feed items periodically, so it's better to have prebuilt URL for that instead of making it each time
	getManyURL := apiURL
	qParams := getManyURL.Query()
	qParams.Add("count", strconv.Itoa(batchSize))
	getManyURL.RawQuery = qParams.Encode()

	return &fetcher{
		apiURL:     apiURL,
		getManyURL: getManyURL,
		cli: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

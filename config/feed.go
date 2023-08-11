package config

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

var (
	ErrNoApiURL         = errors.New("no API url is set")
	ErrInvalidUrlFormat = errors.New("invalid URL format")
)

// the regexp might be inaccurate, but according to tests it does its job
var httpUrlRegExp = regexp.MustCompile(`(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?\/[a-zA-Z0-9]{2,}|((https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?)|(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}(\.[a-zA-Z0-9]{2,})? `)

type Feed struct {
	ApiURL    url.URL
	BatchSize int
}

func (f *Feed) load(envPrefix string) error {
	v := setupViper(envPrefix)

	apiURL := v.GetString("api.url")
	if apiURL == "" {
		return ErrNoApiURL
	}

	if !httpUrlRegExp.MatchString(apiURL) {
		return ErrInvalidUrlFormat
	}

	parsed, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	f.ApiURL = *parsed

	v.SetDefault("batch.size", 50)
	f.BatchSize = v.GetInt("batch.size")

	return nil
}

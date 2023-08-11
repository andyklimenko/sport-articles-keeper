package api

import (
	"net/url"
	"strconv"
)

func parseLimitOffset(queryParams url.Values) (int64, int64, error) {
	var (
		limit, offset int64
		err           error
	)
	if l := queryParams.Get("limit"); l != "" {
		limit, err = strconv.ParseInt(l, 10, 64)
		if err != nil {
			return 0, 0, errParseLimit
		}
	}

	if o := queryParams.Get("offset"); o != "" {
		offset, err = strconv.ParseInt(o, 10, 64)
		if err != nil {
			return 0, 0, errParseOffset
		}
	}

	return limit, offset, nil
}

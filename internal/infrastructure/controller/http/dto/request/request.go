package request

import (
	"net/http"
	"strconv"
)

func ParseListParams(r *http.Request) (int, int, error) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset, nil
}

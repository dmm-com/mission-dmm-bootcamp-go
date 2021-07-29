package request

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// Read path parameter `id`
func IDOf(r *http.Request) (int64, error) {
	ids := chi.URLParam(r, "id")

	if ids == "" {
		return -1, errors.Errorf("id was not presence")
	}

	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return -1, errors.Errorf("id was not number")
	}

	return id, nil
}

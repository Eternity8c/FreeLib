package core_http_utils

import (
	core_errors "FreeLib/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param=`%s` by key=`%s` noy a valid integer: %v: %w",
			param, key, err, core_errors.ErrInvalidArgumment,
		)
	}

	return &val, nil
}

func GetStringQueryParam(r *http.Request, key string) string {
	param := r.URL.Query().Get(key)

	return param
}

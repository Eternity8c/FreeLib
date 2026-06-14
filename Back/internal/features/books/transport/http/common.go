package books_transport_http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	core_jwt "github.com/Eternity8c/FreeLib/internal/core/jwt"
	core_http_utils "github.com/Eternity8c/FreeLib/internal/core/transport/http/utils"
)

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get `limit` query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get `offset` query param: %w", err)
	}

	return limit, offset, nil
}

func getIDQueryParam(r *http.Request) (int, error) {
	id, err := core_http_utils.GetIntQueryParam(r, "ID")
	if err != nil {
		return domain.UninitializedID, fmt.Errorf("get `id` query param: %w", err)
	}

	return *id, nil
}

func idFromJWTToken(ctx context.Context) (int, error) {
	claims, ok := core_jwt.ClaimsFromContext(ctx)
	if !ok {
		return domain.UninitializedID, fmt.Errorf("failed claims from context")
	}
	return claims.ID, nil
}

func getGenreQueryParam(r *http.Request) string {
	genre := core_http_utils.GetStringQueryParam(r, "genre")
	return genre
}

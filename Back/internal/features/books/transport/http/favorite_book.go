package books_transport_http

import (
	"FreeLib/internal/core/domain"
	core_jwt "FreeLib/internal/core/jwt"
	core_logger "FreeLib/internal/core/logger"
	core_http_request "FreeLib/internal/core/transport/http/request"
	core_http_responce "FreeLib/internal/core/transport/http/responce"
	"context"
	"fmt"
	"net/http"
)

type FavoriteBookRequest struct {
	BookID int `json:"book_id"`
}

type FavoriteBookResponce struct {
	UserID int             `json:"user_id"`
	Book   BookDTOResponce `json:"book"`
}

func (h *BooksHTTPHandler) FavoriteBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke favorite book handler")

	userID, err := idFromJWTToken(ctx)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed get id from JWT token")
		return
	}

	var request FavoriteBookRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responceHandler.ErrorResponce(err, "failed decode and validate request")
		return
	}

	uID, bookDomain, err := h.bookServices.FavoriteBook(ctx, userID, request.BookID)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed favorite book")
		return
	}

	responce := FavoriteBookResponce{
		UserID: uID,
		Book:   bookDTOFromDomain(bookDomain),
	}
	responceHandler.JSONResponce(responce, http.StatusOK)
}

func idFromJWTToken(ctx context.Context) (int, error) {
	claims, ok := core_jwt.ClaimsFromContext(ctx)
	if !ok {
		return domain.UninitializedID, fmt.Errorf("failed claims from context")
	}
	return claims.ID, nil
}

package books_transport_http

import (
	core_logger "FreeLib/internal/core/logger"
	core_http_responce "FreeLib/internal/core/transport/http/responce"
	core_http_utils "FreeLib/internal/core/transport/http/utils"
	"fmt"
	"net/http"
)

type GetBooksResponce []BookDTOResponce

func (h *BooksHTTPHandler) GetBooks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke GetBooks handler")

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get limit offset query param")
		return
	}

	booksDomain, err := h.bookServices.GetBooks(ctx, limit, offset)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get books")
		return
	}

	responce := GetBooksResponce(booksDTOFromDomains(booksDomain))
	responceHandler.JSONResponce(responce, http.StatusOK)
}

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

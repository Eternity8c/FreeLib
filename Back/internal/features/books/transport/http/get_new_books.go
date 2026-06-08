package books_transport_http

import (
	core_logger "FreeLib/internal/core/logger"
	core_http_responce "FreeLib/internal/core/transport/http/responce"
	"net/http"
)

type GetNewBooksResponce []BookDTOResponce

func (h *BooksHTTPHandler) GetNewBooks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke GetNewBooks handler")

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get limit offset query param")
		return
	}

	booksDomain, err := h.bookServices.GetNewBooks(ctx, limit, offset)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get new books")
	}

	responce := GetNewBooksResponce(booksDTOFromDomains(booksDomain))
	responceHandler.JSONResponce(responce, http.StatusOK)
}

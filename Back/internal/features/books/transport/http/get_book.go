package books_transport_http

import (
	core_logger "FreeLib/internal/core/logger"
	core_http_responce "FreeLib/internal/core/transport/http/responce"
	"net/http"
)

type GetBookResponce BookDTOResponce

func (h *BooksHTTPHandler) GetBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke get book")
	id, err := getIDQueryParam(r)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get ID query param")
		return
	}

	bookDomain, err := h.bookServices.GetBook(ctx, id)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get book")
		return
	}

	responce := GetBookResponce(bookDTOFromDomain(bookDomain))
	responceHandler.JSONResponce(responce, http.StatusOK)
}

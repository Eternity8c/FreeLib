package books_transport_http

import (
	"FreeLib/internal/core/domain"
	core_logger "FreeLib/internal/core/logger"
	core_http_request "FreeLib/internal/core/transport/http/request"
	core_http_responce "FreeLib/internal/core/transport/http/responce"
	"net/http"
)

type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

type CreateBookResponce BookDTOResponce

func (h *BooksHTTPHandler) CreateBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke CreateBook handler")

	var request CreateBookRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responceHandler.ErrorResponce(err, "failed to validate and decode HTTP request")
		return
	}

	bookDomain, err := h.bookServices.CreateBook(ctx, domainFromDto(request))
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to create book")
		return
	}

	responce := CreateBookResponce(bookDTOFromDomain(bookDomain))
	responceHandler.JSONResponce(responce, http.StatusCreated)
}

func domainFromDto(request CreateBookRequest) domain.Book {
	return domain.NewBookUninitialized(request.Title, request.Author, request.Genre)
}

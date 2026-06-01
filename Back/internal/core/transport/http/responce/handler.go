package core_http_responce

import (
	core_errors "FreeLib/internal/core/errors"
	core_logger "FreeLib/internal/core/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HTTPResponceHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponceHandler(
	log *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponceHandler {
	return &HTTPResponceHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponceHandler) JSONResponce(responceBody any, statusCode int) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responceBody); err != nil {
		h.log.Error("write HTTP responce", zap.Error(err))
	}
}

func (h *HTTPResponceHandler) ErrorResponce(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgumment):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponce(statusCode, err, msg)
}

func (h *HTTPResponceHandler) PanicResponce(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponce(statusCode, err, msg)
}

func (h *HTTPResponceHandler) errorResponce(statusCode int, err error, msg string) {
	responce := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponce(responce, statusCode)
}

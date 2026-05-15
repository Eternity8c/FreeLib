package core_http_responce

import (
	core_logger "FreeLib/internal/core/logger"
	"encoding/json"
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

func (h *HTTPResponceHandler) PanicResponce(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.rw.WriteHeader(statusCode)

	responce := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	if err := json.NewEncoder(h.rw).Encode(responce); err != nil {
		h.log.Error("write HTTP responce", zap.Error(err))
	}
}

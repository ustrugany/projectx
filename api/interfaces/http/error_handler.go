package http

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type errorHandler struct {
	logger zap.SugaredLogger
}

func (h errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
}

func (h errorHandler) serve500Error(err error, content string, w http.ResponseWriter) {
	h.logger.Errorw("server error", "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err = fmt.Fprintln(w, content); err != nil {
		h.logger.Fatalw("failed to sent response", "err", err)
	}
}

func (h errorHandler) serve400Error(err error, content string, w http.ResponseWriter) {
	h.logger.Errorw("bad request", "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if _, err = fmt.Fprintln(w, content); err != nil {
		h.logger.Fatalw("failed to sent response", "err", err)
	}
}

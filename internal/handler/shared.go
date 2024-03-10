package handler

import (
	"log/slog"
	"net/http"
)

func MakeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

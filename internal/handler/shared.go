package handler

import (
	"log/slog"
	"net/http"

	"dreampicai/types"

	"github.com/a-h/templ"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func getAuthenticatedUser(r *http.Request) types.AuthenticatedUser {
	user, ok := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}

func MakeHandler(operation string, h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	return otelhttp.NewHandler(http.HandlerFunc(handler), operation).ServeHTTP
}

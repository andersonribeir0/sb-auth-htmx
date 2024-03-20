package handler

import (
	"net/http"

	"dreampicai/cmd/web/view/auth"
)

func (s *Server) HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Login().Render(r.Context(), w)
}

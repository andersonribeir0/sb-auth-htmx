package handler

import (
	"net/http"

	"dreampicai/cmd/web/view/settings"
)

func (s *Server) HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(r, w, settings.Index(user))
}

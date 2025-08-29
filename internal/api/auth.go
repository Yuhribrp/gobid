package api

import (
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/Yuhribrp/gobid/internal/jsonutils"
)

func (api *Api) HandleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"csrf_token": token,
	})
}

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticatedUserId") {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"error": "unauthorized",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

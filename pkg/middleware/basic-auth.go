package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/service"
)

// NewBasicAuth is a middleware to basic auth HTTP request credentials
func NewBasicAuth(user, pass string) func(inner http.Handler) http.Handler {
	expectedUsernameHash := sha256.Sum256([]byte(user))
	expectedPasswordHash := sha256.Sum256([]byte(pass))
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				jsonErrors(w, "Unauthorized", 401)
				return
			}

			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if !(usernameMatch && passwordMatch) {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				jsonErrors(w, "Unauthorized", 401)
				return
			}

			ctx := r.Context()
			user, err := service.Lookup().GetOrRegisterUser(ctx, user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, constant.ContextUser, *user)
			ctx = context.WithValue(ctx, constant.ContextUserID, *user.ID)
			ctx = context.WithValue(ctx, constant.ContextIsAdmin, true)
			inner.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

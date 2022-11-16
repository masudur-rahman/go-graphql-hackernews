package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/masudur-rahman/hackernews/internal/users"
	"github.com/masudur-rahman/hackernews/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func MiddleWare() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// validate jwt token
			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			user := users.User{Username: username}
			id, err := users.GetUserIDByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user.ID = strconv.Itoa(id)
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. Requires Middleware to be executed first
func ForContext(ctx context.Context) *users.User {
	user, _ := ctx.Value(userCtxKey).(*users.User)
	return user
}

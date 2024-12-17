package gateway

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/GoJobs/utils"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			unauthorizedErrorResponse(w, r, fmt.Errorf("authorization token missing"))
			return
		}

		claims, err := utils.ValidateToken(token)

		if err != nil {
			unauthorizedErrorResponse(w, r, fmt.Errorf("authorization token missing"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "isAdmin", claims.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

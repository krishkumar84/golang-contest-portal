package middleware

import (
	"context"
	"net/http"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/types"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/utils/response"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	StudentIDKey contextKey = "student_id"
	RoleKey      contextKey = "role"      
)

type AuthMiddleware struct {
	secretKey []byte
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: []byte(secretKey),
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
			return
		}

		tokenString := cookie.Value
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return m.secretKey, nil
		})

		if err != nil || !token.Valid {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(err))
			return
		}

		// Update the context values to use the custom keys
		ctx := context.WithValue(r.Context(), UserIDKey, claims["user_id"])
		ctx = context.WithValue(ctx, StudentIDKey, claims["student_id"])
		ctx = context.WithValue(ctx, RoleKey, claims["role"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        role := r.Context().Value(RoleKey)
        if role != string(types.RoleAdmin) {
            response.WriteJson(w, http.StatusForbidden, response.GeneralError(fmt.Errorf("admin access required")))
            return
        }
        next.ServeHTTP(w, r)
    })
}

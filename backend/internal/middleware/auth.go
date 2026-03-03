package middleware

import (
	"context"
	"net/http"
	"strings"
	"swiftly/backend/internal/database"
	"swiftly/backend/internal/pkg/auth"
	"swiftly/backend/internal/pkg/response"
)

type contextKey string

const UserIDKey contextKey = "userID"
const SessionIDKey contextKey = "sessionID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := ""

		// 1. Try to get token from Authorization Header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// 2. If not found in Header, try to get from Cookie
		if tokenString == "" {
			cookie, err := r.Cookie("access_token")
			if err == nil {
				tokenString = cookie.Value
			}
		}

		// 3. If still empty, return error
		if tokenString == "" {
			response.Error(w, http.StatusUnauthorized, "Authentication required", nil)
			return
		}

		// Validate Token signature and claims
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		// SECURITY CHECK: Verify token against Redis
		ctx := r.Context()
		rdb := database.GetRedis()
		
		// 1. Check if specific token is blacklisted (e.g., after manual logout)
		blacklisted, _ := rdb.Exists(ctx, "blacklist:"+tokenString).Result()
		if blacklisted > 0 {
			response.Error(w, http.StatusUnauthorized, "Session expired. Please login again.", nil)
			return
		}

		// 2. Check if the parent session still exists (e.g., after remote revocation)
		sessionKey := "session:" + claims.UserID + ":" + claims.SessionID
		exists, _ := rdb.Exists(ctx, sessionKey).Result()
		if exists == 0 {
			response.Error(w, http.StatusUnauthorized, "Session has been revoked", nil)
			return
		}

		ctx = context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, SessionIDKey, claims.SessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

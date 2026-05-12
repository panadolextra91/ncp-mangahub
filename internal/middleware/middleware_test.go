package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/internal/middleware"
	"github.com/user/mangahub/pkg/auth"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret"
	mw := middleware.AuthMiddleware(secret)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(middleware.RoleKey).(string)
		assert.True(t, ok)
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		assert.True(t, ok)
		assert.Equal(t, "admin", role)
		assert.Equal(t, 1, userID)
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Valid Token", func(t *testing.T) {
		tokenString, _ := auth.GenerateToken(1, "tester", "admin", secret)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		rr := httptest.NewRecorder()

		mw(nextHandler).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Missing Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		mw(nextHandler).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid")
		rr := httptest.NewRecorder()
		mw(nextHandler).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

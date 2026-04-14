package auth_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/pkg/auth"
)

func TestJWT(t *testing.T) {
	secret := "test-secret"
	userID := 123
	role := "admin"

	t.Run("Generate and Validate", func(t *testing.T) {
		token, err := auth.GenerateToken(userID, role, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := auth.ValidateToken(token, secret)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("Invalid Secret", func(t *testing.T) {
		token, _ := auth.GenerateToken(userID, role, secret)
		claims, err := auth.ValidateToken(token, "wrong-secret")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Invalid Token Format", func(t *testing.T) {
		claims, err := auth.ValidateToken("not-a-token", secret)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

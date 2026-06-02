package core_jwt

import (
	"FreeLib/internal/core/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user domain.User) (string, error) {
	config := NewConfigMust()

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(config.Expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}
	return signedToken, nil
}

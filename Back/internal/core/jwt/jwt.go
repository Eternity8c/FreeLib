package core_jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/Eternity8c/FreeLib/internal/core/domain"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID      int  `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

type contextKey string

const claimsContextKey contextKey = "core_jwt_claims"

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsContextKey).(*Claims)
	return c, ok
}

func GenerateToken(user domain.User) (string, error) {
	config := NewConfigMust()
	now := time.Now()

	claims := &Claims{
		ID:      user.ID,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(config.Expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}
	return signedToken, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	config := NewConfigMust()
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

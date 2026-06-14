package core_http_middleware

import (
	"net/http"
	"strings"

	core_errors "github.com/Eternity8c/FreeLib/internal/core/errors"
	core_jwt "github.com/Eternity8c/FreeLib/internal/core/jwt"
	core_logger "github.com/Eternity8c/FreeLib/internal/core/logger"
	core_http_responce "github.com/Eternity8c/FreeLib/internal/core/transport/http/responce"
)

func Authorization() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responceHandler := core_http_responce.NewHTTPResponceHandler(log, w)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responceHandler.ErrorResponce(core_errors.ErrUnauthorized, "user unauthorized")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := core_jwt.ParseToken(tokenString)
			if err != nil {
				responceHandler.ErrorResponce(core_errors.ErrUnauthorized, "invalid token")
				return
			}

			ctx = core_jwt.ContextWithClaims(ctx, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := core_logger.FromContext(ctx)
		responceHandler := core_http_responce.NewHTTPResponceHandler(log, w)
		claims, ok := core_jwt.ClaimsFromContext(ctx)
		if !ok {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responceHandler.ErrorResponce(core_errors.ErrUnauthorized, "user unauthorized")
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			var err error
			claims, err = core_jwt.ParseToken(tokenString)
			if err != nil {
				responceHandler.ErrorResponce(core_errors.ErrUnauthorized, "invalid token")
				return
			}
		}

		if !claims.IsAdmin {
			responceHandler.ErrorResponce(core_errors.ErrForbidden, "access denied: admins only")
			return
		}

		next(w, r)
	}
}

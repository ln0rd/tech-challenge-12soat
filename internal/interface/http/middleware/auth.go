package middleware

import (
	"context"
	"net/http"
	"strings"

	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	tokenService domain.TokenService
	logger       *zap.Logger
}

func NewAuthMiddleware(tokenService domain.TokenService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		logger:       logger,
	}
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		am.logger.Info("Authenticating request", zap.String("path", r.URL.Path))

		// Extrai o token do header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			am.logger.Error("Missing Authorization header")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Verifica se o header começa com "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			am.logger.Error("Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Extrai o token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			am.logger.Error("Empty token")
			http.Error(w, "Empty token", http.StatusUnauthorized)
			return
		}

		// Valida o token
		claims, err := am.tokenService.ValidateToken(token)
		if err != nil {
			am.logger.Error("Invalid token", zap.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		am.logger.Info("Token validated successfully", zap.String("email", claims.Email))

		// Adiciona as claims ao contexto da requisição
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)

		// Chama o próximo handler com o contexto atualizado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

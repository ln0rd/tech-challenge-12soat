package middleware

import (
	"net/http"

	"github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
	"go.uber.org/zap"
)

type AuthorizationMiddleware struct {
	logger *zap.Logger
}

func NewAuthorizationMiddleware(logger *zap.Logger) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		logger: logger,
	}
}

// Middleware para verificar se o usuário é mechanic (acesso total)
func (am *AuthorizationMiddleware) RequireMechanic(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*auth.Claims)
		if !ok {
			am.logger.Error("Claims not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.UserType != "mechanic" {
			am.logger.Error("User is not mechanic",
				zap.String("userType", claims.UserType),
				zap.String("userID", claims.UserID.String()))
			http.Error(w, "Forbidden: Only mechanics can access this resource", http.StatusForbidden)
			return
		}

		am.logger.Info("Mechanic access granted",
			zap.String("userType", claims.UserType),
			zap.String("userID", claims.UserID.String()),
			zap.String("path", r.URL.Path))

		next.ServeHTTP(w, r)
	}
}

// Middleware para verificar se o usuário é vehicle_owner (acesso limitado)
func (am *AuthorizationMiddleware) RequireVehicleOwner(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*auth.Claims)
		if !ok {
			am.logger.Error("Claims not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.UserType != "vehicle_owner" {
			am.logger.Error("User is not vehicle_owner",
				zap.String("userType", claims.UserType),
				zap.String("userID", claims.UserID.String()))
			http.Error(w, "Forbidden: Only vehicle owners can access this resource", http.StatusForbidden)
			return
		}

		am.logger.Info("Vehicle owner access granted",
			zap.String("userType", claims.UserType),
			zap.String("userID", claims.UserID.String()),
			zap.String("path", r.URL.Path))

		next.ServeHTTP(w, r)
	}
}

// Middleware para verificar se o usuário é admin (acesso total)
func (am *AuthorizationMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*auth.Claims)
		if !ok {
			am.logger.Error("Claims not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.UserType != "admin" {
			am.logger.Error("User is not admin",
				zap.String("userType", claims.UserType),
				zap.String("userID", claims.UserID.String()))
			http.Error(w, "Forbidden: Only admins can access this resource", http.StatusForbidden)
			return
		}

		am.logger.Info("Admin access granted",
			zap.String("userType", claims.UserType),
			zap.String("userID", claims.UserID.String()),
			zap.String("path", r.URL.Path))

		next.ServeHTTP(w, r)
	}
}

// Middleware para verificar se o usuário é mechanic OU admin
func (am *AuthorizationMiddleware) RequireMechanicOrAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*auth.Claims)
		if !ok {
			am.logger.Error("Claims not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.UserType != "mechanic" && claims.UserType != "admin" {
			am.logger.Error("User is not mechanic or admin",
				zap.String("userType", claims.UserType),
				zap.String("userID", claims.UserID.String()))
			http.Error(w, "Forbidden: Only mechanics and admins can access this resource", http.StatusForbidden)
			return
		}

		am.logger.Info("Mechanic or admin access granted",
			zap.String("userType", claims.UserType),
			zap.String("userID", claims.UserID.String()),
			zap.String("path", r.URL.Path))

		next.ServeHTTP(w, r)
	}
}

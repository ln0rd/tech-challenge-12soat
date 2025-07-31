package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/auth"
	"go.uber.org/zap"
)

type JWTService struct {
	secretKey     []byte
	refreshSecret []byte
	logger        *zap.Logger
}

func NewJWTService(logger *zap.Logger) *JWTService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production"
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = "your-refresh-secret-key-change-in-production"
	}

	return &JWTService{
		secretKey:     []byte(secretKey),
		refreshSecret: []byte(refreshSecret),
		logger:        logger,
	}
}

func (j *JWTService) GenerateToken(userInfo domain.UserInfo) (string, error) {
	j.logger.Info("Generating JWT token", zap.String("email", userInfo.Email))

	now := time.Now()
	exp := now.Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id":   userInfo.ID.String(),
		"email":     userInfo.Email,
		"username":  userInfo.Username,
		"user_type": userInfo.UserType,
		"exp":       exp.Unix(),
		"iat":       now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		j.logger.Error("Error signing token", zap.Error(err))
		return "", err
	}

	j.logger.Info("JWT token generated successfully")
	return tokenString, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*domain.Claims, error) {
	j.logger.Info("Validating JWT token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		j.logger.Error("Error parsing token", zap.Error(err))
		return nil, err
	}

	if !token.Valid {
		j.logger.Error("Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		j.logger.Error("Invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		j.logger.Error("Invalid user ID in token", zap.Error(err))
		return nil, err
	}

	domainClaims := &domain.Claims{
		UserID:   userID,
		Email:    claims["email"].(string),
		Username: claims["username"].(string),
		UserType: claims["user_type"].(string),
		Exp:      int64(claims["exp"].(float64)),
		Iat:      int64(claims["iat"].(float64)),
	}

	j.logger.Info("JWT token validated successfully", zap.String("email", domainClaims.Email))
	return domainClaims, nil
}

func (j *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	j.logger.Info("Generating refresh token", zap.String("userID", userID.String()))

	// Gera um refresh token simples (em produção, use uma abordagem mais segura)
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		j.logger.Error("Error generating refresh token", zap.Error(err))
		return "", err
	}

	tokenString := base64.URLEncoding.EncodeToString(refreshToken)
	j.logger.Info("Refresh token generated successfully")
	return tokenString, nil
}

func (j *JWTService) ValidateRefreshToken(refreshToken string) (uuid.UUID, error) {
	j.logger.Info("Validating refresh token")

	// Em uma implementação real, você validaria o refresh token contra o banco de dados
	// Por simplicidade, aqui apenas decodificamos
	_, err := base64.URLEncoding.DecodeString(refreshToken)
	if err != nil {
		j.logger.Error("Invalid refresh token", zap.Error(err))
		return uuid.Nil, fmt.Errorf("invalid refresh token")
	}

	// Em uma implementação real, você extrairia o userID do refresh token
	// Por simplicidade, retornamos um UUID vazio
	j.logger.Info("Refresh token validated successfully")
	return uuid.Nil, nil
}

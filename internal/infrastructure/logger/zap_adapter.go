package logger

import (
	domain "github.com/ln0rd/tech_challenge_12soat/internal/domain/interfaces"
	"go.uber.org/zap"
)

// ZapAdapter adapta o zap.Logger para a interface Logger do domínio
type ZapAdapter struct {
	logger *zap.Logger
}

// NewZapAdapter cria uma nova instância do adaptador
func NewZapAdapter(logger *zap.Logger) domain.Logger {
	return &ZapAdapter{
		logger: logger,
	}
}

// Info implementa o método Info da interface Logger
func (z *ZapAdapter) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

// Error implementa o método Error da interface Logger
func (z *ZapAdapter) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

// Warn implementa o método Warn da interface Logger
func (z *ZapAdapter) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

// Debug implementa o método Debug da interface Logger
func (z *ZapAdapter) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

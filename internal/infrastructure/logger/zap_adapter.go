package logger

import (
	"go.uber.org/zap"
)

// Logger define a interface para logging
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
}

// ZapAdapter implementa Logger usando zap.Logger
type ZapAdapter struct {
	logger *zap.Logger
}

// NewZapAdapter cria uma nova instância do adaptador
func NewZapAdapter(logger *zap.Logger) Logger {
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

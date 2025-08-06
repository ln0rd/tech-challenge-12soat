package mocks

import (
	"go.uber.org/zap"
)

// LoggerMock implementa Logger para testes
type LoggerMock struct {
	InfoFunc  func(msg string, fields ...zap.Field)
	ErrorFunc func(msg string, fields ...zap.Field)
	WarnFunc  func(msg string, fields ...zap.Field)
	DebugFunc func(msg string, fields ...zap.Field)
}

// Info chama a função mock
func (m *LoggerMock) Info(msg string, fields ...zap.Field) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg, fields...)
	}
}

// Error chama a função mock
func (m *LoggerMock) Error(msg string, fields ...zap.Field) {
	if m.ErrorFunc != nil {
		m.ErrorFunc(msg, fields...)
	}
}

// Warn chama a função mock
func (m *LoggerMock) Warn(msg string, fields ...zap.Field) {
	if m.WarnFunc != nil {
		m.WarnFunc(msg, fields...)
	}
}

// Debug chama a função mock
func (m *LoggerMock) Debug(msg string, fields ...zap.Field) {
	if m.DebugFunc != nil {
		m.DebugFunc(msg, fields...)
	}
}

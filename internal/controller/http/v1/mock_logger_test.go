package v1

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// implements the logger.Interface
type MockLogger struct {
	mock.Mock
	T *testing.T
}

func NewMockLogger(t *testing.T) *MockLogger {
	return &MockLogger{T: t}
}

func (m *MockLogger) Debug(message interface{}, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Info(message string, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Warn(message string, args ...interface{}) {
	m.Called(message, args)
}

func (m *MockLogger) Error(message interface{}, args ...interface{}) {
	m.Called(message, args)

	if m.T != nil {
		if err, ok := message.(error); ok {
			m.T.Logf("Error logged: %v, args: %v", err, args)
		} else {
			m.T.Logf("Error logged: %v, args: %v", message, args)
		}
	}
}

func (m *MockLogger) Fatal(message interface{}, args ...interface{}) {
	m.Called(message, args)
}

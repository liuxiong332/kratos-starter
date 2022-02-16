package logger

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestLogger(t *testing.T) {
	logger := NewLogger()
	logger.Log(log.LevelInfo, "hello", "Hello world")
}

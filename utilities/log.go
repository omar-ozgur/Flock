package utilities

import (
	"go.uber.org/zap"
)

// Log utilities
var Logger, _ = zap.NewDevelopment()
var Sugar = Logger.Sugar()

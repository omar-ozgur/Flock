package utilities

import (
	"go.uber.org/zap"
)

// Logger is a lightweight logger
var Logger, _ = zap.NewDevelopment()

// Sugar is a logger that wraps Logger to provide a more 'ergonomic' API
var Sugar = Logger.Sugar()

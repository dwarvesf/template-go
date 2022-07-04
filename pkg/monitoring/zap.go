package monitoring

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newZapLogger returns zapLogger struct based on default config
func newZapLogger(debugMode bool) *zap.Logger {
	return zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(newEncoderConfig()),
		os.Stdout,
		newLoggingLevel(debugMode),
	))
}

// newEncoderConfig returns opinionated EncoderConfig
func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:  "message",                   // to show {"message": "content"}
		LevelKey:    "level",                     // to show {"level": "info"}
		TimeKey:     "time",                      // to show {"time":"2021-06-21T09:25:51+08:00"}
		LineEnding:  zapcore.DefaultLineEnding,   // character to separate multi lines
		EncodeLevel: zapcore.CapitalLevelEncoder, // to show {"level": "info"} or to show {"level": "INFO"}
		EncodeTime:  zapcore.ISO8601TimeEncoder,  // to format time as {"timestamp":"2021-06-21T09:25:51.230+08:00"}
	}
}

func newLoggingLevel(debugMode bool) zapcore.Level {
	if debugMode {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

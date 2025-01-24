package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	*os.File
}

func NewLoggerZap(file *os.File) *Logger {
	writeSyncer := zapcore.AddSync(file)
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)

	return &Logger{zap.New(core), file}

}

func (l *Logger) Close() {
	l.Logger.Sync()
	l.File.Close()
}

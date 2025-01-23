package zaplog

import (
	"fmt"
	"messenger/config"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Глобальная переменная для логгера
var Logger *zap.Logger

// Инициализация логгера
func Init() {
	logFile := fmt.Sprintf("%s/app_%s.log", config.LogDir, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Настроим zapcore для вывода логов в файл
	writeSyncer := zapcore.AddSync(file)
	encoderConfig := zap.NewProductionEncoderConfig()

	// Добавляем дату, время и другие данные
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Создаем новый core с уровнем логирования и выводом в файл
	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)

	// Создаем логгер с указанной конфигурацией
	Logger = zap.New(core)

	// Логируем, что инициализация прошла успешно
	Logger.Info("Logger инициализирован успешно")
}

func LogHandler(key string, val string) zap.Field {
	return zap.String(key, val)
}

func LogError(err error) zap.Field {
	return zap.Error(err)
}

func LogString(key string, val string) zap.Field {
	return zap.String(key, val)
}

// Закрытие логгера
func Shutdown() {
	if Logger != nil {
		Logger.Sync()
	}
}

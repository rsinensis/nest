package logger

import (
	"io"
	"os"

	"github.com/rsinensis/nest/module/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// InitLogger init logger
func InitLogger(mode string) {
	cfg := setting.GetSetting()

	filename := cfg.Section("log").Key("Filename").MustString("log/app.log")
	maxSize := cfg.Section("log").Key("MaxSize").MustInt(100)
	maxBackups := cfg.Section("log").Key("MaxBackups").MustInt(15)
	maxAge := cfg.Section("log").Key("MaxAge").MustInt(28)
	level := cfg.Section("log").Key("Level").MustString("info")
	var encoder zapcore.EncoderConfig
	var output io.Writer
	switch mode {
	case "dev":
		encoder = zap.NewDevelopmentEncoderConfig()
		output = os.Stdout
	case "prod":
	case "test":
		encoder = zap.NewProductionEncoderConfig()
		output = &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize, // megabytes
			MaxBackups: maxBackups,
			MaxAge:     maxAge, // days
		}
	}

	w := zapcore.AddSync(output)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		w,
		getZapLevel(level),
	)
	logger = zap.New(core)
}

// getZapLevel change string to zapcore.Level
func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug", "DEBUG":
		return zap.DebugLevel
	case "info", "INFO", "": // make the zero value useful
		return zap.InfoLevel
	case "warn", "WARN":
		return zap.WarnLevel
	case "error", "ERROR":
		return zap.ErrorLevel
	case "dpanic", "DPANIC":
		return zap.DPanicLevel
	case "panic", "PANIC":
		return zap.PanicLevel
	case "fatal", "FATAL":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// GetLogger return logger
func GetLogger() *zap.Logger{
	return logger
}
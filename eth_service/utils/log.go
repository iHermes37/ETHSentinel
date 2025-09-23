package utils

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "sync"
)

var (
    log  *zap.Logger
    once sync.Once
)

// Init 初始化日志系统，只会执行一次
func Init(level string) *zap.Logger {
    once.Do(func() {
        var zapLevel zapcore.Level
        switch level {
        case "debug":
            zapLevel = zap.DebugLevel
        case "info":
            zapLevel = zap.InfoLevel
        case "warn":
            zapLevel = zap.WarnLevel
        case "error":
            zapLevel = zap.ErrorLevel
        default:
            zapLevel = zap.InfoLevel
        }

        config := zap.Config{
            Level:       zap.NewAtomicLevelAt(zapLevel),
            Development: false,
            Encoding:    "json", // 可改为 "console"
            OutputPaths: []string{"stdout", "logs/app.log"}, // 控制台 + 文件
            ErrorOutputPaths: []string{"stderr", "logs/error.log"},
            EncoderConfig: zapcore.EncoderConfig{
                TimeKey:        "time",
                LevelKey:       "level",
                NameKey:        "logger",
                CallerKey:      "caller",
                MessageKey:     "msg",
                StacktraceKey:  "stacktrace",
                LineEnding:     zapcore.DefaultLineEnding,
                EncodeLevel:    zapcore.CapitalLevelEncoder,
                EncodeTime:     zapcore.ISO8601TimeEncoder,
                EncodeDuration: zapcore.StringDurationEncoder,
                EncodeCaller:   zapcore.ShortCallerEncoder,
            },
        }

        var err error
        log, err = config.Build()
        if err != nil {
            panic(err)
        }
    })
    return log
}

// Get 获取全局日志
func Getlog() *zap.Logger {
    if log == nil {
        return Init("info")
    }
    return log
}

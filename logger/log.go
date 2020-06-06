package logger

import (
	"os"
	"log"
	"helloweb/common"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Z global zap.Logger
var Z *zap.Logger

const loggerTimeLayout = "2006-01-02 15:04:05.000"

// logPath log file path
var logPath string

// NewEncoderConfig  新编码器配置
func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// TimeEncoder 时间编码器
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(loggerTimeLayout))
}

// InitLogger initialize zap logger
func InitLogger() {
	if Z != nil {
		return
	}
	if common.HomePath == ""{
		log.Fatal("初始化app路径失败，HomePath = ",common.HomePath)
	}
	logPath = common.HomePath + "/" + "log.log"

	hook := lumberjack.Logger{
		Filename:   logPath, // 日志文件路径
		MaxSize:    128,     // 128MB
		MaxBackups: 30,      // 最多保留30个备份
		MaxAge:     7,       // 文件存储时长，7天
		Compress:   true,    // 是否压缩
	}

	w := zapcore.AddSync(&hook)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),
			w),
		zap.DebugLevel,
	)

	Z = zap.New(core, zap.AddCaller())
}

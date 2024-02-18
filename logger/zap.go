package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
	"strings"
)

var ZapLogger *zap.Logger

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志文件夹
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名

	MaxAge       int  `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                      // 日志留存时间
	ShowLine     bool `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole bool `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
}

func (z *Zap) ZapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (z *Zap) TransportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func InitZapLogger() {
	z := &Zap{
		Level:         "debug",
		Prefix:        "",
		Format:        "json",
		Director:      "log",
		EncodeLevel:   "LowercaseLevelEncoder",
		StacktraceKey: "stacktrace",
		MaxAge:        3,
		ShowLine:      true,
		LogInConsole:  true,
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  z.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    z.ZapEncodeLevel(),
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	var level zapcore.Level
	if z.LogInConsole {
		level = zapcore.DebugLevel
	} else {
		level = z.TransportLevel()
	}

	var cores []zapcore.Core
	if z.ShowLine {
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}
	consoleDebugging := zapcore.Lock(os.Stdout)
	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleDebugging, level))

	fileWriter := zapcore.AddSync(os.Stdout)

	if runtime.GOOS == "windows" {
		fileWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:   path.Join(z.Director, "zap.log"),
			MaxSize:    1 << 30,
			MaxBackups: 3,
			MaxAge:     z.MaxAge,
			Compress:   true,
		})
	} else {
		fileWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:   path.Join(z.Director, "zap.log"),
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     z.MaxAge,
			Compress:   true,
		})
	}
	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), fileWriter, level))

	core := zapcore.NewTee(cores...)
	ZapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	ZapLogger.Info("log 初始化成功")
}

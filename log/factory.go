package log

import (
	"errors"
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultPrefix       = "default"                                    // 默认文件名前缀
	defaultMaxAge       = "168h"                                       // 默认文件最大保存时间：7天
	defaultRotationTime = "24h"                                        // 默认文件切割时间间隔：24小时
	defaultConfigs      = []*Config{{Type: "console", Level: "debug"}} // 默认配置，仅在终端输出Debug以上等级的日志
)

// Config 日志配置
type Config struct {
	Type  string `yaml:"type"`  // 日志类型：console/file
	Level string `yaml:"level"` // 日志等级：debug/info/warn/error/fatal

	// 以下配置仅在类型为file时生效
	Prefix       string `yaml:"prefix"`        // 文件名前缀，例如prefix为tmp，则文件名为tmp.log
	MaxAge       string `yaml:"max-age"`       // 文件最大保存时间，使用time.ParseDuration函数进行计算
	RotationTime string `yaml:"rotation-time"` // 文件切割时间间隔，使用time.ParseDuration函数进行计算
	RotationSize int64  `yaml:"rotation-size"` // 文件最大大小，单位MB
}

// New 创建新的日志器
func New(configs ...*Config) (Logger, error) {
	return NewWithCallerSkip(1, configs...)
}

// NewDefaultLogger 使用新的配置覆盖默认日志器
func NewDefaultLogger(configs ...*Config) error {
	logger, err := NewWithCallerSkip(2, configs...)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

// NewWithCallerSkip 创建新的日志器，额外加上skip调用者的个数
func NewWithCallerSkip(skip int, configs ...*Config) (Logger, error) {
	if configs == nil || len(configs) == 0 {
		configs = defaultConfigs
	}

	hasConsole := false
	cores := make([]zapcore.Core, 0)
	for _, cnf := range configs {
		if cnf.Level == "" {
			cnf.Level = "info"
		}
		lvl, err := zapcore.ParseLevel(cnf.Level)
		if err != nil {
			return nil, fmt.Errorf("parse log level failed: %w", err)
		}

		switch cnf.Type {
		case "console":
			if !hasConsole {
				ecnf := zap.NewProductionEncoderConfig()
				ecnf.EncodeTime = zapcore.ISO8601TimeEncoder
				ecnf.EncodeLevel = zapcore.CapitalColorLevelEncoder
				enc := zapcore.NewConsoleEncoder(ecnf)

				core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), lvl)
				cores = append(cores, core)
				hasConsole = true
			}
		case "file":
			options := make([]rotatelogs.Option, 0)
			// TODO prefix重复时的操作
			if cnf.Prefix == "" {
				cnf.Prefix = defaultPrefix
			}
			filename := cnf.Prefix + ".log"
			options = append(options, rotatelogs.WithLinkName(filename))

			if cnf.MaxAge == "" {
				cnf.MaxAge = defaultMaxAge
			}
			maxAge, err := time.ParseDuration(cnf.MaxAge)
			if err != nil {
				return nil, fmt.Errorf("parse log max age failed: %s", cnf.MaxAge)
			}
			options = append(options, rotatelogs.WithMaxAge(maxAge))

			if cnf.RotationTime == "" {
				cnf.RotationTime = defaultRotationTime
			}
			rotationTime, err := time.ParseDuration(cnf.RotationTime)
			if err != nil {
				return nil, fmt.Errorf("parse rotation time failed: %s", cnf.RotationTime)
			}
			options = append(options, rotatelogs.WithRotationTime(rotationTime))

			if cnf.RotationSize > 0 {
				options = append(options, rotatelogs.WithRotationSize(cnf.RotationSize*1024*1024))
			}

			writer, err := rotatelogs.New(filename+".%Y%m%d%H%M", options...)
			if err != nil {
				return nil, fmt.Errorf("create log writer failed: %w", err)
			}

			// 创建日志编码器
			ecnf := zap.NewProductionEncoderConfig()
			ecnf.EncodeTime = zapcore.ISO8601TimeEncoder
			ecnf.EncodeLevel = zapcore.CapitalLevelEncoder
			enc := zapcore.NewConsoleEncoder(ecnf)
			// 创建日志核心
			core := zapcore.NewCore(enc, zapcore.AddSync(writer), lvl)
			cores = append(cores, core)
		default:
			return nil, errors.New("unexpected log type: " + cnf.Type)
		}
	}

	core := zapcore.NewTee(cores...)
	return &zapLoggerWrapper{zap.New(core, zap.AddCaller(), zap.AddCallerSkip(skip))}, nil
}

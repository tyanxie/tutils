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

const (
	defaultPrefix       = "log"  // 默认文件名前缀
	defaultMaxAge       = "168h" // 默认文件最大保存时间：7天
	defaultRotationTime = "24h"  // 默认文件切割时间间隔：24小时
)

var defaultConfigs = []*Config{{Type: "console", Level: "debug"}} // 默认配置，仅在终端输出Debug以上等级的日志

// Config 日志配置
type Config struct {
	Type  string // 日志类型：console/file
	Level string // 日志等级：debug/info/warn/error/fatal

	// 以下配置仅在类型为file时生效
	Prefix       string // 文件名前缀，例如prefix为tmp，则文件名为tmp.log
	MaxAge       string // 文件最大保存时间，使用time.ParseDuration函数进行计算
	RotationTime string // 文件切割时间间隔，使用time.ParseDuration函数进行计算
	RotationSize int64  // 文件最大大小，单位MB
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
	prefixes := make(map[string]struct{}, 0)
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
			} else {
				return nil, errors.New("duplicate console writer")
			}
		case "file":
			options := make([]rotatelogs.Option, 0)
			// 默认前缀
			if cnf.Prefix == "" {
				cnf.Prefix = defaultPrefix
			}
			// 校验文件前缀
			if _, ok := prefixes[cnf.Prefix]; ok {
				return nil, fmt.Errorf("duplicate file prefix: %s", cnf.Prefix)
			}
			prefixes[cnf.Prefix] = struct{}{}
			// 文件名称
			filename := cnf.Prefix + ".log"
			options = append(options, rotatelogs.WithLinkName(filename))
			// 最大保存时间
			if cnf.MaxAge == "" {
				cnf.MaxAge = defaultMaxAge
			}
			maxAge, err := time.ParseDuration(cnf.MaxAge)
			if err != nil {
				return nil, fmt.Errorf("parse log max age failed: %s", cnf.MaxAge)
			}
			options = append(options, rotatelogs.WithMaxAge(maxAge))
			// 文件切割时间间隔
			if cnf.RotationTime == "" {
				cnf.RotationTime = defaultRotationTime
			}
			rotationTime, err := time.ParseDuration(cnf.RotationTime)
			if err != nil {
				return nil, fmt.Errorf("parse rotation time failed: %s", cnf.RotationTime)
			}
			options = append(options, rotatelogs.WithRotationTime(rotationTime))
			// 文件最大大小
			if cnf.RotationSize > 0 {
				options = append(options, rotatelogs.WithRotationSize(cnf.RotationSize*1024*1024))
			}
			// 文件滚动
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
	// 构造日志器
	core := zapcore.NewTee(cores...)
	return &zapLoggerWrapper{zap.New(core, zap.AddCaller(), zap.AddCallerSkip(skip))}, nil
}

package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	OutputTypeConsole = "console" // 日志输出类型：console
	OutputTypeFile    = "file"    // 日志输出类型：文件
)

const (
	defaultFilename   = "output.log" // 默认文件名称
	defaultMaxSize    = 300          // 默认文件最大大小：300MB
	defaultMaxAge     = 7            // 默认文件最大保存时间：7天
	defaultMaxBackups = 10           // 默认文件最大保存数量：10个
)

var (
	defaultConfigs     = []*Config{{Type: "console", Level: "debug"}}           // 默认配置，仅在终端输出Debug以上等级的日志
	defaultTimeEncoder = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") // 默认时间格式化
)

// Config 日志配置
type Config struct {
	Type  string // 日志类型：console/file
	Level string // 日志等级：debug/info/warn/error/fatal

	// 以下配置仅在类型为file时生效，核心逻辑为按照文件大小切割文件
	Path       string // 文件目录路径
	Filename   string // 文件名称
	MaxSize    int64  // 文件最大大小，单位MB
	MaxAge     int64  // 文件最大保存时间，单位天
	MaxBackups int64  // 文件最大保存数量
	Compress   bool   // 是否压缩文件
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
	filenames := make(map[string]struct{})
	for _, cnf := range configs {
		if cnf.Level == "" {
			cnf.Level = "info"
		}
		lvl, err := zapcore.ParseLevel(cnf.Level)
		if err != nil {
			return nil, fmt.Errorf("parse log level failed: %w", err)
		}

		switch cnf.Type {
		case OutputTypeConsole:
			if !hasConsole {
				// 创建日志编码器
				ecnf := zap.NewProductionEncoderConfig()
				ecnf.EncodeTime = defaultTimeEncoder
				ecnf.EncodeLevel = zapcore.CapitalColorLevelEncoder
				enc := zapcore.NewConsoleEncoder(ecnf)
				// 创建日志核心
				core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), lvl)
				cores = append(cores, core)
				hasConsole = true
			} else {
				return nil, errors.New("duplicate console writer")
			}
		case OutputTypeFile:
			// 文件名称
			if cnf.Filename == "" {
				cnf.Filename = defaultFilename
			}
			if cnf.Path != "" {
				cnf.Filename = filepath.Join(cnf.Path, cnf.Filename)
			}
			// 校验文件名
			if _, ok := filenames[cnf.Filename]; ok {
				return nil, fmt.Errorf("duplicate filename: %s", cnf.Filename)
			}
			filenames[cnf.Filename] = struct{}{}
			// 默认文件最大大小
			if cnf.MaxSize <= 0 {
				cnf.MaxSize = defaultMaxSize
			}
			// 默认最大保存时间
			if cnf.MaxAge <= 0 {
				cnf.MaxAge = defaultMaxAge
			}
			// 默认文件最大保存数量
			if cnf.MaxBackups <= 0 {
				cnf.MaxBackups = defaultMaxBackups
			}
			// 创建文件
			writer, err := NewWriter(cnf.Filename, cnf.MaxSize, cnf.MaxAge, cnf.MaxBackups, cnf.Compress)
			if err != nil {
				return nil, fmt.Errorf("create writer failed: %w", err)
			}
			// 创建日志编码器
			ecnf := zap.NewProductionEncoderConfig()
			ecnf.EncodeTime = defaultTimeEncoder
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

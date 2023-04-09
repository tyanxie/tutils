package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

var defaultLogger Logger // 默认日志器

type contextKey struct{} // 日志器在context中的key

// init 初始化创建默认日志器，失败时panic
func init() {
	var err error
	defaultLogger, err = NewWithCallerSkip(2, defaultConfigs...)
	if err != nil {
		panic(fmt.Sprintf("init tlog failed: %+v", err))
	}
}

// Debug 打印Debug等级的日志
func Debug(msg string) {
	defaultLogger.Debug(msg)
}

// DebugContext 打印Debug等级的日志，优先使用ctx中的logger
func DebugContext(ctx context.Context, msg string) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Debug(msg)
}

// Info 打印Info等级的日志
func Info(msg string) {
	defaultLogger.Info(msg)
}

// InfoContext 打印Info等级的日志，优先使用ctx中的logger
func InfoContext(ctx context.Context, msg string) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Info(msg)
}

// Warn 打印Warn等级的日志
func Warn(msg string) {
	defaultLogger.Warn(msg)
}

// WarnContext 打印Warn等级的日志，优先使用ctx中的logger
func WarnContext(ctx context.Context, msg string) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Warn(msg)
}

// Error 打印Error等级的日志
func Error(msg string) {
	defaultLogger.Error(msg)
}

// ErrorContext 打印Error等级的日志，优先使用ctx中的logger
func ErrorContext(ctx context.Context, msg string) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Error(msg)
}

// Fatal 打印Fatal等级的日志
func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

// FatalContext 打印Fatal等级的日志，优先使用ctx中的logger
func FatalContext(ctx context.Context, msg string) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Fatal(msg)
}

// Debugf 使用fmt.Sprintf打印Debug等级的日志
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// DebugContextf 使用fmt.Sprintf打印Debug等级的日志，优先使用ctx中的logger
func DebugContextf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Debugf(format, args...)
}

// Infof 使用fmt.Sprintf打印Info等级的日志
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// InfoContextf 使用fmt.Sprintf打印Info等级的日志，优先使用ctx中的logger
func InfoContextf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Infof(format, args...)
}

// Warnf 使用fmt.Sprintf打印Warn等级的日志
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// WarnContextf 使用fmt.Sprintf打印Warn等级的日志，优先使用ctx中的logger
func WarnContextf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Warnf(format, args...)
}

// Errorf 使用fmt.Sprintf打印Error等级的日志
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// ErrorContextf 使用fmt.Sprintf打印Error等级的日志，优先使用ctx中的logger
func ErrorContextf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Errorf(format, args...)
}

// Fatalf 使用fmt.Sprintf打印Fatal等级的日志
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// FatalContextf 使用fmt.Sprintf打印Fatal等级的日志，优先使用ctx中的logger
func FatalContextf(ctx context.Context, format string, args ...interface{}) {
	logger := FromContext(ctx)
	if l, ok := logger.(*zapLoggerWrapper); ok {
		logger = &zapLoggerWrapper{
			l.Logger.WithOptions(zap.AddCallerSkip(1)),
		}
	}
	logger.Fatalf(format, args...)
}

// WithField 向日志增加一个自定义字段
func WithField(key string, value interface{}) Logger {
	return &zapLoggerWrapper{
		defaultLogger.(*zapLoggerWrapper).Logger.
			WithOptions(zap.AddCallerSkip(-1)).
			With(zap.Any(key, value)),
	}
}

// WithContextField 向ctx中的日志器中增加一个自定义字段，不会修改ctx中的日志器
func WithContextField(ctx context.Context, key string, value interface{}) Logger {
	return FromContext(ctx).WithField(key, value)
}

// WithFields 向日志增加多个自定义字段
func WithFields(fields Fields) Logger {
	fs := make([]zap.Field, len(fields))
	i := 0
	for key, value := range fields {
		fs[i] = zap.Any(key, value)
		i++
	}
	return &zapLoggerWrapper{
		defaultLogger.(*zapLoggerWrapper).Logger.
			WithOptions(zap.AddCallerSkip(-1)).
			With(fs...),
	}
}

// WithContextFields 向ctx中的日志器中增加多个自定义字段，不会修改ctx中的日志器
func WithContextFields(ctx context.Context, fields Fields) Logger {
	return FromContext(ctx).WithFields(fields)
}

// ToContext 向ctx中放入日志器
func ToContext(ctx context.Context, logger Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, contextKey{}, logger)
	return ctx
}

// FromContext 从ctx中获取日志器，如果没有则返回默认日志器
func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return &zapLoggerWrapper{
			defaultLogger.(*zapLoggerWrapper).Logger.WithOptions(zap.AddCallerSkip(-1)),
		}
	}
	logger, ok := ctx.Value(contextKey{}).(Logger)
	if ok {
		return logger
	}
	return &zapLoggerWrapper{
		defaultLogger.(*zapLoggerWrapper).Logger.WithOptions(zap.AddCallerSkip(-1)),
	}
}

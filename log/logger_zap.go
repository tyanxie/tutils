package log

import "go.uber.org/zap"

// zapLoggerWrapper 包装了zap日志的日志器
type zapLoggerWrapper struct {
	*zap.Logger
}

// Debug 打印Debug等级的日志
func (logger *zapLoggerWrapper) Debug(msg string) {
	logger.Logger.Debug(msg)
}

// Info 打印Info等级的日志
func (logger *zapLoggerWrapper) Info(msg string) {
	logger.Logger.Info(msg)
}

// Warn 打印Warn等级的日志
func (logger *zapLoggerWrapper) Warn(msg string) {
	logger.Logger.Warn(msg)
}

// Error 打印Error等级的日志
func (logger *zapLoggerWrapper) Error(msg string) {
	logger.Logger.Error(msg)
}

// Fatal 打印Fatal等级的日志
func (logger *zapLoggerWrapper) Fatal(msg string) {
	logger.Logger.Fatal(msg)
}

// Debugf 使用fmt.Sprintf打印Debug等级的日志
func (logger *zapLoggerWrapper) Debugf(template string, args ...interface{}) {
	logger.Logger.Sugar().Debugf(template, args...)
}

// Infof 使用fmt.Sprintf打印Info等级的日志
func (logger *zapLoggerWrapper) Infof(template string, args ...interface{}) {
	logger.Logger.Sugar().Infof(template, args...)
}

// Warnf 使用fmt.Sprintf打印Warn等级的日志
func (logger *zapLoggerWrapper) Warnf(template string, args ...interface{}) {
	logger.Logger.Sugar().Warnf(template, args...)
}

// Errorf 使用fmt.Sprintf打印Error等级的日志
func (logger *zapLoggerWrapper) Errorf(template string, args ...interface{}) {
	logger.Logger.Sugar().Errorf(template, args...)
}

// Fatalf 使用fmt.Sprintf打印Fatal等级的日志
func (logger *zapLoggerWrapper) Fatalf(template string, args ...interface{}) {
	logger.Logger.Sugar().Fatalf(template, args...)
}

// WithField 向日志增加一个自定义字段
func (logger *zapLoggerWrapper) WithField(key string, value interface{}) Logger {
	return &zapLoggerWrapper{
		logger.Logger.With(zap.Any(key, value)),
	}
}

// WithFields 向日志增加多个自定义字段
func (logger *zapLoggerWrapper) WithFields(fields Fields) Logger {
	if fields == nil || len(fields) == 0 {
		return logger
	}
	fs := make([]zap.Field, len(fields))
	i := 0
	for key, value := range fields {
		fs[i] = zap.Any(key, value)
		i++
	}
	return &zapLoggerWrapper{
		logger.Logger.With(fs...),
	}
}

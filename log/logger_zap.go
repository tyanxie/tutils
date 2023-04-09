package log

import "go.uber.org/zap"

// zapLoggerWrapper 包装了zap日志的日志器
type zapLoggerWrapper struct {
	*zap.Logger
}

func (logger *zapLoggerWrapper) Debug(msg string) {
	logger.Logger.Debug(msg)
}

func (logger *zapLoggerWrapper) Info(msg string) {
	logger.Logger.Info(msg)
}

func (logger *zapLoggerWrapper) Warn(msg string) {
	logger.Logger.Warn(msg)
}

func (logger *zapLoggerWrapper) Error(msg string) {
	logger.Logger.Error(msg)
}

func (logger *zapLoggerWrapper) Fatal(msg string) {
	logger.Logger.Fatal(msg)
}

func (logger *zapLoggerWrapper) Debugf(template string, args ...interface{}) {
	logger.Logger.Sugar().Debugf(template, args...)
}

func (logger *zapLoggerWrapper) Infof(template string, args ...interface{}) {
	logger.Logger.Sugar().Infof(template, args...)
}

func (logger *zapLoggerWrapper) Warnf(template string, args ...interface{}) {
	logger.Logger.Sugar().Warnf(template, args...)
}

func (logger *zapLoggerWrapper) Errorf(template string, args ...interface{}) {
	logger.Logger.Sugar().Errorf(template, args...)
}

func (logger *zapLoggerWrapper) Fatalf(template string, args ...interface{}) {
	logger.Logger.Sugar().Fatalf(template, args...)
}

func (logger *zapLoggerWrapper) WithField(key string, value interface{}) Logger {
	return &zapLoggerWrapper{
		logger.Logger.With(zap.Any(key, value)),
	}
}

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

func (logger *zapLoggerWrapper) WithError(err error) Logger {
	return &zapLoggerWrapper{
		logger.Logger.With(zap.Error(err)),
	}
}

func (logger *zapLoggerWrapper) Named(name string) Logger {
	if name == "" {
		return logger
	}
	return &zapLoggerWrapper{
		logger.Logger.Named(name),
	}
}

package log

// Fields 用户定义的字段
type Fields map[string]interface{}

// Logger 日志接口
type Logger interface {
	// Debug 打印Debug等级的日志
	Debug(msg string)
	// Info 打印Info等级的日志
	Info(msg string)
	// Warn 打印Warn等级的日志
	Warn(msg string)
	// Error 打印Error等级的日志
	Error(msg string)
	// Fatal 打印Fatal等级的日志
	Fatal(msg string)

	// Debugf 使用fmt.Sprintf打印Debug等级的日志
	Debugf(template string, args ...interface{})
	// Infof 使用fmt.Sprintf打印Info等级的日志
	Infof(template string, args ...interface{})
	// Warnf 使用fmt.Sprintf打印Warn等级的日志
	Warnf(template string, args ...interface{})
	// Errorf 使用fmt.Sprintf打印Error等级的日志
	Errorf(template string, args ...interface{})
	// Fatalf 使用fmt.Sprintf打印Fatal等级的日志
	Fatalf(template string, args ...interface{})

	// WithField 向日志增加一个自定义字段
	WithField(key string, value interface{}) Logger
	// WithFields 向日志增加多个自定义字段
	WithFields(fields Fields) Logger
	// WithError 向日志增加error错误类型字段
	WithError(err error) Logger
	// Named 向日志器增加标题
	Named(name string) Logger
}

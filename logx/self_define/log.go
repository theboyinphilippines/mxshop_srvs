package self_define

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// 抽象工程设计模式，自定义logger接口，与zap解耦，可以替换其他日志框架
type Logger interface {
	Debug(msg string)
	DebugC(context context.Context, msg string)
	Debugf(format string, args ...interface{})
	DebugW(msg string, keyAndValues ...interface{})
	DebugWC(context context.Context, msg string, keyAndValues ...interface{})
}

var _ Logger = &zapLogger{}

// 使用zaplogger实例来实现接口
type zapLogger struct {
	zapLogger *zap.Logger
}

func (z *zapLogger) Debug(msg string) {
	z.zapLogger.Debug(msg)
}

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func (z *zapLogger) DebugC(context context.Context, msg string) {
	panic("implement me")
}

func (z *zapLogger) Debugf(format string, args ...interface{}) {
	panic("implement me")
}

func (z *zapLogger) DebugW(msg string, keyAndValues ...interface{}) {
	panic("implement me")
}

func (z *zapLogger) DebugWC(context context.Context, msg string, keyAndValues ...interface{}) {
	panic("implement me")
}

var (
	defaultLogger = New(NewOptions())
	mu            sync.Mutex
)

// 根据传入的options，new一个logger实例
func New(opt *Options) *zapLogger {
	if opt == nil {
		//没有传入选项，就用options的默认选项
		opt = NewOptions()
	}
	//实例化zap
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opt.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	loggerConfig := zap.Config{
		Level: zap.NewAtomicLevelAt(zapLevel),
	}
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		panic(any(err))
	}
	return &zapLogger{
		zapLogger: l.Named(opt.Name), //可以传入logger的名称
	}
}

// 获取初始化的logger实例，供调用
func GetLogger() *zapLogger {
	return defaultLogger

}
func Init(opt *Options) {
	mu.Lock()
	defer mu.Unlock()
	defaultLogger = New(opt)
}

package self_define

import (
	"fmt"
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	FORMAT_CONSOLE = "console"
	FORMAT_JSON    = "json"
	OUTPUT_STD     = "stdout"
	OUTPUT_STD_ERR = "stderr"
	flagLevel      = "log.level"
)

// log的配置文件选项
type Options struct {
	OutputPaths    []string `mapstructure:"output_paths" json:"output_paths"`
	ErrOutputPaths []string `mapstructure:"err_output_paths" json:"err_output_paths"`
	Level          string   `mapstructure:"level" json:"level"`
	Format         string   `mapstructure:"format" json:"format"`
	Name           string   `mapstructure:"name" json:"name"`
}

type Option func(o *Options)

func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// 使用函数选项模式，设置options
func NewOptions(opts ...Option) *Options {
	//默认设置
	options := &Options{
		Level:          zapcore.InfoLevel.String(),
		Format:         FORMAT_CONSOLE,
		OutputPaths:    []string{OUTPUT_STD},
		ErrOutputPaths: []string{OUTPUT_STD_ERR},
	}
	//循环传入选项参数，改变设置
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// 校验Options的各个字段
func (o *Options) Validate() []error {
	var errs []error
	format := strings.ToLower(o.Format)
	if format != FORMAT_CONSOLE && format != FORMAT_JSON {
		errs = append(errs, fmt.Errorf("not support format %s", o.Format))
	}
	return errs
}

// Options的各个字段,可以从命令行传入，列映射到flag字段上
func (o *Options) AddFlags(fs pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "log level")
}

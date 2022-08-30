/**
 * @Author raven
 * @Description
 * @Date 2022/8/29
 **/
package log

import (
	"context"
)

var defaultLogger = DaenerysLogger{}

func Debug(ctx context.Context, msg string) {
	defaultLogger.Debug(ctx, msg)
}
func Debugf(ctx context.Context, format string, arg ...interface{}) {
	defaultLogger.Debugf(ctx, format, arg)
}

func Info(ctx context.Context, msg string) {
	defaultLogger.Info(ctx, msg)
}

func Infof(ctx context.Context, format string, arg ...interface{}) {
	defaultLogger.Infof(ctx, format, arg)
}

func Warn(ctx context.Context, msg string) {
	defaultLogger.Warn(ctx, msg)
}
func Warnf(ctx context.Context, format string, arg ...interface{}) {
	defaultLogger.Warnf(ctx, format, arg)
}

func Error(ctx context.Context, msg string) {
	defaultLogger.Error(ctx, msg)
}
func Errorf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Errorf(ctx, format, args...)
}

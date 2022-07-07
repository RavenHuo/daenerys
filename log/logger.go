/**
 * @Author raven
 * @Description
 * @Date 2022/7/7
 **/
package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

func Error(ctx context.Context, msg string) {
	logrus.Error(msg)
}
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
func Info(ctx context.Context, msg string) {
	logrus.Info(msg)
}
func Infof(ctx context.Context, format string, arg ...interface{}) {
	logrus.Infof(format, arg)
}
func Warn(ctx context.Context, msg string) {
	logrus.Warn(msg)
}
func Warnf(ctx context.Context, format string, arg ...interface{}) {
	logrus.Warnf(format, arg)
}

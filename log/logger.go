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

const TraceIdField = "trace-id"

func Error(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Error(msg)
}
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Errorf(format, args...)
}
func Info(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Infof(msg)
}
func Infof(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Infof(format, arg)
}
func Warn(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Warn(msg)
}
func Warnf(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Warnf(format, arg)
}

func getTraceId(ctx context.Context) string {
	traceId := ctx.Value(TraceIdField)
	if traceId == nil {
		traceId = ""
	}
	return traceId.(string)
}

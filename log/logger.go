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

type DaenerysLogger struct{}

func (d DaenerysLogger) Error(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Error(msg)
}
func (d DaenerysLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Errorf(format, args...)
}
func (d DaenerysLogger) Info(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Info(msg)
}
func (d DaenerysLogger) Infof(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Infof(format, arg...)
}
func (d DaenerysLogger) Warn(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Warn(msg)
}
func (d DaenerysLogger) Warnf(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Warnf(format, arg...)
}

func (d DaenerysLogger) Debug(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Debug(msg)
}
func (d DaenerysLogger) Debugf(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).Debugf(format, arg...)
}

func getTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ctx.Value(TraceIdField)
	if traceId == nil {
		traceId = ""
	}
	return traceId.(string)
}

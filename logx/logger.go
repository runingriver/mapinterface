package logx

import "context"

var (
	logIns = &Logger{}
)

type Logger struct {
	logger LogItf
}

func SetLogger(log LogItf) *Logger {
	logIns.logger = log
	return logIns
}

func Fatal(format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.Fatal(format, v...)
}

func Error(format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.Error(format, v...)
}
func Warn(format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.Warn(format, v...)
}
func Notice(format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.Notice(format, v...)
}
func Info(format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.Info(format, v...)
}
func Debug(format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.Debug(format, v...)
}
func Trace(format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.Trace(format, v...)
}

func CtxFatal(ctx context.Context, format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.CtxFatal(ctx, format, v...)
}
func CtxError(ctx context.Context, format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.CtxError(ctx, format, v...)
}
func CtxWarn(ctx context.Context, format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.CtxWarn(ctx, format, v...)
}
func CtxNotice(ctx context.Context, format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.CtxNotice(ctx, format, v...)
}
func CtxInfo(ctx context.Context, format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.CtxInfo(ctx, format, v...)
}
func CtxDebug(ctx context.Context, format string, v ...interface{}) {
	if logIns.logger == nil {
		return
	}
	logIns.logger.CtxDebug(ctx, format, v...)
}

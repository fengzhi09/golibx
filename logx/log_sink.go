package logx

//import (
//	"context"
//	"time"
//)
//
//type ILogSink[T int | string] interface {
//	Debugf(ctx context.Context, format string, args ...any)
//	Infof(ctx context.Context, format string, args ...any)
//	Warnf(ctx context.Context, format string, args ...any)
//	Errorf(ctx context.Context, format string, args ...any)
//	Panicf(ctx context.Context, format string, args ...any)
//	SetLevel(level T)
//	//WithField(key string, value any)ILogSink[T int | string]
//}
//
//type ILogSinkCloser[T int | string] interface {
//	ILogSink[T]
//	Close(ctx context.Context) error
//}
//type logProxy[T int | string] struct {
//	sink   ILogSinkCloser[T]
//	lvlmap map[T]LogLevel
//	lvl    LogLevel
//}
//
//func (l logProxy[T]) SetLevel(level T) {
//	l.sink.SetLevel(level)
//	l.lvl = l.lvlmap[level]
//}
//
//func (l logProxy[T]) SetLevelLogx(level LogLevel) {
//	l.lvl = level
//	for k, v := range l.lvlmap {
//		if v == level {
//			l.sink.SetLevel(k)
//			break
//		}
//	}
//}
//
//func (l logProxy[T]) Write(event *LogEvent) error {
//
//	if event.Level >= l.lvl {
//		switch event.Level {
//		case DEBUG:
//			l.sink.Debugf(event.Ctx, event.Msg)
//		case INFO:
//			l.sink.Infof(event.Ctx, event.Msg)
//		case WARN:
//			l.sink.Warnf(event.Ctx, event.Msg)
//		case ERROR:
//			l.sink.Errorf(event.Ctx, event.Msg)
//		case PANIC:
//			l.sink.Panicf(event.Ctx, event.Msg)
//		}
//	}
//	return nil
//}
//
//func SinkWriter[T int | string](sink ILogSink[T], lvlmap map[T]LogLevel) LogWriter {
//	return &logProxy[T]{
//		sink:   sink,
//		lvlmap: lvlmap,
//	}
//}
//func (l logProxy[T]) Debugf(ctx context.Context, format string, args ...any) {
//	l.Write(&LogEvent{
//		Ctx:    ctx,
//		Level:  DEBUG,
//		Msg:    format,
//		Time:   time.Now(),
//		Caller: CallerOfName("logProxy"),
//		With:   nil,
//	})
//}
//
//func (l logProxy[T]) Infof(ctx context.Context, format string, args ...any) {
//	l.sink.Infof(ctx, format, args...)
//}
//
//func (l logProxy[T]) Warnf(ctx context.Context, format string, args ...any) {
//	l.sink.Warnf(ctx, format, args...)
//}
//
//func (l logProxy[T]) Errorf(ctx context.Context, format string, args ...any) {
//	l.sink.Errorf(ctx, format, args...)
//}
//
//func (l logProxy[T]) Panicf(ctx context.Context, format string, args ...any) {
//	l.sink.Panicf(ctx, format, args...)
//}
//
//func (l logProxy[T]) Close(ctx context.Context) error {
//	if closer, ok := l.sink.(ILogSinkCloser[T]); ok {
//		return closer.Close(ctx)
//	}
//	return nil
//}

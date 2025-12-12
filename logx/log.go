package logx

import (
	"context"
	"fmt"
	"github.com/fengzhi09/golibx/gox"
	"os"
	"sync"
	"time"
)

// Log 日志工具类
type Log struct {
	sync.Mutex
	writers   []LogWriter
	level     LogLevel
	moduleMap map[string]LogLevel
	With      map[string]any
}

var (
	instance *Log
	once     sync.Once
)

// Init 初始化日志
func (l *Log) Init(app string, cmdLvl LogLevel, rotates ...*RotateConf) error {
	l.Lock()
	defer l.Unlock()
	l.level = cmdLvl
	l.moduleMap = make(map[string]LogLevel)
	l.writers = make([]LogWriter, 0)
	for _, rotate := range rotates {
		if rotate != nil && rotate.Enable {
			l.writers = append(l.writers, RotateWriter(app, rotate))
		}
	}
	if len(l.writers) == 0 {
		l.writers = append(l.writers, CmdWriter(cmdLvl, true))
	}
	return nil
}

// log 写入日志
func (l *Log) log(ctx context.Context, level LogLevel, module, message string) {
	if int(level) < int(l.level) {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	caller := gox.CallerSkip("logx")
	event := NewEvent(level, ctx, message)
	event.With = map[string]any{"module": module}
	for k, v := range l.With {
		event.With[k] = v
	}
	locate := fmt.Sprintf("%s,at %s,%d`", module, caller.File, caller.Line)
	logLine := fmt.Sprintf(" [%s] [%s] [%s] %s\n", timestamp, level, locate, message)
	for _, writer := range l.writers {
		err := writer.Write(event)
		if err != nil {
			// 如果写入失败，尝试写入标准错误
			_, _ = fmt.Fprintf(os.Stderr, "日志写入失败: %v, 原日志内容: %s\n", err, logLine)
		}
	}
}

// Debug 调试日志
func (l *Log) Debugf(ctx context.Context, format string, args ...any) {
	l.log(ctx, DEBUG, "ROOT", fmt.Sprintf(format, args...))
}

// Info 信息日志
func (l *Log) Infof(ctx context.Context, format string, args ...any) {
	l.log(ctx, INFO, "ROOT", fmt.Sprintf(format, args...))
}

// Warn 警告日志
func (l *Log) Warnf(ctx context.Context, format string, args ...any) {
	l.log(ctx, WARN, "ROOT", fmt.Sprintf(format, args...))
}

// Error 错误日志
func (l *Log) Errorf(ctx context.Context, format string, args ...any) {
	l.log(ctx, ERROR, "ROOT", fmt.Sprintf(format, args...))
}

// Debug 调试日志
func (l *Log) Debug(ctx context.Context, format string, args ...any) {
	l.log(ctx, DEBUG, "ROOT", fmt.Sprintf(format, args...))
}

// Info 信息日志
func (l *Log) Info(ctx context.Context, message string) {
	l.log(ctx, INFO, "ROOT", message)
}

// Warn 警告日志
func (l *Log) Warn(ctx context.Context, message string) {
	l.log(ctx, WARN, "ROOT", message)
}

// Error 错误日志
func (l *Log) Error(ctx context.Context, message string) {
	l.log(ctx, ERROR, "ROOT", message)
}

// Debug 模块调试日志
func (l *Log) DebugM(ctx context.Context, module string, message string) {
	l.log(ctx, DEBUG, module, message)
}

// InfoM 模块信息日志
func (l *Log) InfoM(ctx context.Context, module, message string) {
	l.log(ctx, INFO, module, message)
}

// WarnM 模块警告日志
func (l *Log) WarnM(ctx context.Context, module, message string) {
	l.log(ctx, WARN, module, message)
}

// ErrorM 模块错误日志
func (l *Log) ErrorM(ctx context.Context, module, message string) {
	l.log(ctx, ERROR, module, message)
}

// Debug 模块调试日志
func (l *Log) DebugfM(ctx context.Context, module string, format string, args ...any) {
	l.log(ctx, DEBUG, module, fmt.Sprintf(format, args...))
}

// InfofM 模块信息日志
func (l *Log) InfofM(ctx context.Context, module string, format string, args ...any) {
	l.log(ctx, INFO, module, fmt.Sprintf(format, args...))
}

// WarnfM 模块警告日志
func (l *Log) WarnfM(ctx context.Context, module string, format string, args ...any) {
	l.log(ctx, WARN, module, fmt.Sprintf(format, args...))
}

// ErrorfM 模块错误日志
func (l *Log) ErrorfM(ctx context.Context, module string, format string, args ...any) {
	l.log(ctx, ERROR, module, fmt.Sprintf(format, args...))
}

// CloseLogs 关闭日志
func (l *Log) CloseWriters(ctx context.Context) {
	for _, writer := range l.writers {
		_ = writer.Close(ctx)
	}
	l.writers = nil
	instance = nil
}
func (l *Log) WithField(key string, value any) {
	l.With = map[string]any{
		key: value,
	}
}

// log 获取日志实例
func log() *Log {
	once.Do(func() {
		instance = &Log{
			writers:   []LogWriter{},
			level:     INFO,
			moduleMap: make(map[string]LogLevel),
		}
	})
	return instance
}

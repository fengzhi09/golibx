package logx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fengzhi09/golibx/gox"
)

// LogLevel 日志级别
type LogLevel int

const (
	// DEBUG 调试级别
	DEBUG LogLevel = iota
	// INFO 信息级别
	INFO
	// WARN 警告级别
	WARN
	// ERROR 错误级别
	ERROR
	// Panic 紧急级别
	PANIC
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBU"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERRO"
	case PANIC:
		return "PANI"
	default:
		return "UNKN"
	}
}

func ParseLogLevel(level string) LogLevel {
	level = gox.SubStr(level, 0, 4)
	// 设置日志级别
	switch strings.ToUpper(level) {
	case "DEBU", "TRAC":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERRO":
		return ERROR
	case "PANI", "FATA":
		return PANIC
	default:
		return INFO
	}
}

type LogEvent struct {
	Ctx    context.Context
	Level  LogLevel
	Msg    string
	Time   time.Time
	With   map[string]any
	Caller *gox.CallerInfo
}
type LogWriter interface {
	SetLevelLogx(level LogLevel)
	Write(event *LogEvent) error
	Close(ctx context.Context) error
}

func NewEvent(level LogLevel, ctx context.Context, msg string) *LogEvent {
	return &LogEvent{
		Ctx:    ctx,
		Level:  level,
		Msg:    msg,
		Time:   time.Now(),
		Caller: gox.CallerSkip("logx"),
		With:   nil,
	}
}

// Init 初始化全局日志
func Init(app string, level LogLevel, rotate *RotateConf) error {
	return log().Init(app, level, rotate)
}

func CloseLogs(ctx context.Context) {
	log().CloseWriters(ctx)
}

// Debugf 调试日志
func Debugf(ctx context.Context, format string, args ...any) {
	log().log(ctx, DEBUG, "ROOT", fmt.Sprintf(format, args...))
}

// Infof 信息日志
func Infof(ctx context.Context, format string, args ...any) {
	log().log(ctx, INFO, "ROOT", fmt.Sprintf(format, args...))
}

// Warnf 警告日志
func Warnf(ctx context.Context, format string, args ...any) {
	log().log(ctx, WARN, "ROOT", fmt.Sprintf(format, args...))
}

// Errorf 错误日志
func Errorf(ctx context.Context, format string, args ...any) {
	log().log(ctx, ERROR, "ROOT", fmt.Sprintf(format, args...))
}

// Debug 调试日志
func Debug(ctx context.Context, format string, args ...any) {
	log().log(ctx, DEBUG, "ROOT", fmt.Sprintf(format, args...))
}

// Info 信息日志
func Info(ctx context.Context, message string) {
	log().log(ctx, INFO, "ROOT", message)
}

// Warn 警告日志
func Warn(ctx context.Context, message string) {
	log().log(ctx, WARN, "ROOT", message)
}

// Error 错误日志
func Error(ctx context.Context, message string) {
	log().log(ctx, ERROR, "ROOT", message)
}

// Debug 模块调试日志
func DebugM(ctx context.Context, module string, message string) {
	log().log(ctx, DEBUG, module, message)
}

// InfoM 模块信息日志
func InfoM(ctx context.Context, module, message string) {
	log().log(ctx, INFO, module, message)
}

// WarnM 模块警告日志
func WarnM(ctx context.Context, module, message string) {
	log().log(ctx, WARN, module, message)
}

// ErrorM 模块错误日志
func ErrorM(ctx context.Context, module, message string) {
	log().log(ctx, ERROR, module, message)
}

// Debug 模块调试日志
func DebugfM(ctx context.Context, module string, format string, args ...any) {
	log().log(ctx, DEBUG, module, fmt.Sprintf(format, args...))
}

// InfofM 模块信息日志
func InfofM(ctx context.Context, module string, format string, args ...any) {
	log().log(ctx, INFO, module, fmt.Sprintf(format, args...))
}

// WarnfM 模块警告日志
func WarnfM(ctx context.Context, module string, format string, args ...any) {
	log().log(ctx, WARN, module, fmt.Sprintf(format, args...))
}

// ErrorfM 模块错误日志
func ErrorfM(ctx context.Context, module string, format string, args ...any) {
	log().log(ctx, ERROR, module, fmt.Sprintf(format, args...))
}

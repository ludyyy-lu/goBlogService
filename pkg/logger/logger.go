package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"runtime"
	"time"
)

type Level int8
type Fields map[string]any

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

// 日志分级
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	fields    Fields
	callers   []string
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	l := log.New(w, prefix, flag)
	return &Logger{newLogger: l}
}

func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

// 设置日志公共字段
func (l *Logger) WithFields(f Fields) *Logger {
	l1 := l.clone()
	if l1.fields == nil {
		l1.fields = make(Fields)
	}
	// for k,v := range f {
	// 	l1.fields[k] = v
	// }
	maps.Copy(l1.fields, f)

	return l1
}

// 设置日志上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
	l1 := l.clone()
	l1.ctx = ctx
	return l1
}

// 设置当前某一层调用栈的信息（程序计数器、文件信息和信号）
func (l *Logger) WithCaller(skip int) *Logger {
	l1 := l.clone()
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		l1.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
	}
	return l1
}

// 设置当前的整个调用栈信息
func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		s := fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, s)
		if !more {
			break
		}
	}
	l1 := l.clone()
	l1.callers = callers
	return l1
}

func (l *Logger) JSONFormat(level Level, message string) map[string]any {
	data := make(Fields, len(l.fields)+4)
	data["level"] = level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers

	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

func (l *Logger) Output(level Level, message string) {
	body, _ := json.Marshal(l.JSONFormat(level, message))
	content := string(body)
	switch level {
	case LevelDebug:
		l.newLogger.Print(content)
	case LevelInfo:
		l.newLogger.Print(content)
	case LevelWarn:
		l.newLogger.Print(content)
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Fatal(content)
	case LevelPanic:
		l.newLogger.Panic(content)
	}
}

// debug
func (l *Logger) Debug(v ...any) {
	l.Output(LevelDebug, fmt.Sprint(v...))
}
func (l *Logger) Debugf(format string, v ...any) {
	l.Output(LevelDebug, fmt.Sprintf(format, v...))
}

// info
func (l *Logger) Info(v ...any) {
	l.Output(LevelInfo, fmt.Sprint(v...))
}
func (l *Logger) Infof(format string, v ...any) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

// warn
func (l *Logger) Warn(v ...any) {
	l.Output(LevelWarn, fmt.Sprint(v...))
}
func (l *Logger) Warnf(format string, v ...any) {
	l.Output(LevelWarn, fmt.Sprintf(format, v...))
}

// error
func (l *Logger) Error(v ...any) {
	l.Output(LevelError, fmt.Sprint(v...))
}
func (l *Logger) Errorf(format string, v ...any) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}

// fatal
func (l *Logger) Fatal(v ...any) {
	l.Output(LevelFatal, fmt.Sprint(v...))
}
func (l *Logger) Fatalf(format string, v ...any) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
}

// Panic
func (l *Logger) Panic(v ...any) {
	l.Output(LevelPanic, fmt.Sprint(v...))
}
func (l *Logger) Panicf(format string, v ...any) {
	l.Output(LevelPanic, fmt.Sprintf(format, v...))
}

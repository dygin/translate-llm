package utils

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelFatal
)

// LogEntry 日志条目
type LogEntry struct {
	Level     LogLevel   `json:"level"`
	Timestamp time.Time  `json:"timestamp"`
	Message   string     `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     error      `json:"error,omitempty"`
}

// Log 记录日志
func Log(ctx context.Context, level LogLevel, message string, data interface{}, err error) {
	entry := LogEntry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Data:      data,
		Error:     err,
	}

	// 根据日志级别记录
	switch level {
	case LogLevelDebug:
		g.Log().Debug(ctx, formatLogEntry(entry))
	case LogLevelInfo:
		g.Log().Info(ctx, formatLogEntry(entry))
	case LogLevelWarning:
		g.Log().Warning(ctx, formatLogEntry(entry))
	case LogLevelError:
		g.Log().Error(ctx, formatLogEntry(entry))
	case LogLevelFatal:
		g.Log().Fatal(ctx, formatLogEntry(entry))
	}
}

// Debug 记录调试日志
func Debug(ctx context.Context, message string, data interface{}) {
	Log(ctx, LogLevelDebug, message, data, nil)
}

// Info 记录信息日志
func Info(ctx context.Context, message string, data interface{}) {
	Log(ctx, LogLevelInfo, message, data, nil)
}

// Warning 记录警告日志
func Warning(ctx context.Context, message string, data interface{}) {
	Log(ctx, LogLevelWarning, message, data, nil)
}

// Error 记录错误日志
func Error(ctx context.Context, message string, err error, data interface{}) {
	Log(ctx, LogLevelError, message, data, err)
}

// Fatal 记录致命错误日志
func Fatal(ctx context.Context, message string, err error, data interface{}) {
	Log(ctx, LogLevelFatal, message, data, err)
}

// formatLogEntry 格式化日志条目
func formatLogEntry(entry LogEntry) string {
	levelStr := getLogLevelString(entry.Level)
	timeStr := entry.Timestamp.Format("2006-01-02 15:04:05.000")
	
	msg := fmt.Sprintf("[%s] %s %s", levelStr, timeStr, entry.Message)
	
	if entry.Data != nil {
		msg += fmt.Sprintf(" Data: %v", entry.Data)
	}
	
	if entry.Error != nil {
		msg += fmt.Sprintf(" Error: %v", entry.Error)
	}
	
	return msg
}

// getLogLevelString 获取日志级别字符串
func getLogLevelString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogRequest 记录请求日志
func LogRequest(ctx context.Context, r *ghttp.Request) {
	Info(ctx, "收到请求", map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"ip":         r.GetClientIp(),
		"user_agent": r.UserAgent(),
	})
}

// LogResponse 记录响应日志
func LogResponse(ctx context.Context, r *ghttp.Request, duration time.Duration) {
	Info(ctx, "请求完成", map[string]interface{}{
		"method":     r.Method,
		"path":       r.URL.Path,
		"status":     r.Response.Status,
		"duration":   duration.String(),
	})
}

// LogError 记录错误日志
func LogError(ctx context.Context, r *ghttp.Request, err error) {
	Error(ctx, "请求处理失败", err, map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"ip":     r.GetClientIp(),
	})
} 
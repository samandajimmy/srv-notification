package logger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbs-go/nlogger"
	"io"
	stdLog "log"
	"os"
	"time"
)

type Json struct {
	level     nlogger.LogLevel
	levelStr  map[nlogger.LogLevel]string
	skipTrace int
	writer    *stdLog.Logger
	ioWriter  io.Writer
	namespace string
	flags     int
	ctx       context.Context
}

func (l Json) Fatal(msg string, args ...interface{}) {
	l.print(nlogger.LevelFatal, msg, evaluateOptions(args))
}

func (l Json) Fatalf(format string, args ...interface{}) {
	l.print(nlogger.LevelFatal, format, &Options{
		fmtArgs: args,
	})
}

func (l Json) Error(msg string, args ...interface{}) {
	l.print(nlogger.LevelError, msg, evaluateOptions(args))
}

func (l Json) Errorf(format string, args ...interface{}) {
	l.print(nlogger.LevelError, format, &Options{
		fmtArgs: args,
	})
}

func (l Json) Warn(msg string, args ...interface{}) {
	l.print(nlogger.LevelWarn, msg, evaluateOptions(args))
}

func (l Json) Warnf(format string, args ...interface{}) {
	l.print(nlogger.LevelWarn, format, &Options{
		fmtArgs: args,
	})
}

func (l Json) Info(msg string, args ...interface{}) {
	l.print(nlogger.LevelInfo, msg, evaluateOptions(args))
}

func (l Json) Infof(format string, args ...interface{}) {
	l.print(nlogger.LevelInfo, format, &Options{
		fmtArgs: args,
	})
}

func (l *Json) Debug(msg string, args ...interface{}) {
	l.print(nlogger.LevelDebug, msg, evaluateOptions(args))
}

func (l *Json) Debugf(format string, args ...interface{}) {
	l.print(nlogger.LevelDebug, format, &Options{
		fmtArgs: args,
	})
}

func (l *Json) NewChild(args ...interface{}) nlogger.Logger {
	options := evaluateOptions(args)

	// Override namespace if option is set
	var n string
	if options.namespace != "" {
		n = options.namespace
	} else {
		n = l.namespace
	}

	// Init logger
	logger := NewJson(l.level, l.ioWriter, n)

	// Set context if available
	if options.context != nil {
		logger.ctx = options.context
	} else if l.ctx != nil {
		// Fallback context to logger if set
		logger.ctx = l.ctx
	}

	return logger
}

func NewJson(level nlogger.LogLevel, w io.Writer, namespace string) *Json {
	// If writer is nil, set default writer to Stdout
	if w == nil {
		w = os.Stdout
	}

	// Init standard logger instance
	l := Json{
		level: level,
		levelStr: map[nlogger.LogLevel]string{
			nlogger.LevelFatal: "FATAL",
			nlogger.LevelError: "ERROR",
			nlogger.LevelWarn:  "WARN",
			nlogger.LevelInfo:  "INFO",
			nlogger.LevelDebug: "DEBUG",
		},
		skipTrace: 2,
		writer:    stdLog.New(w, "", 0),
		ioWriter:  w,
		namespace: namespace,
	}
	return &l
}

func (l *Json) print(outLevel nlogger.LogLevel, msg string, options *Options) {
	// if output level is greater than log level, don't print
	if outLevel > l.level {
		return
	}

	// Init json body
	jsonBody := JsonBody{
		Timestamp: time.Now().Format(time.RFC3339),
		LevelId:   outLevel,
		Level:     l.levelStr[outLevel],
	}

	// Compose message
	if len(options.fmtArgs) > 0 {
		jsonBody.Message = fmt.Sprintf(msg, options.fmtArgs...)
	} else {
		jsonBody.Message = msg
	}

	// Set namespace
	if options.namespace != "" {
		jsonBody.Namespace = options.namespace
	}

	// If error exists, then print error and its trace
	if options.err != nil && outLevel <= nlogger.LevelError {
		// Trace caller
		filePath, line := nlogger.Trace(l.skipTrace)

		// Set error
		jsonBody.Error = options.err
		jsonBody.Trace = fmt.Sprintf("%s:%d", filePath, line)

		// Print cause
		if errCause := errors.Unwrap(options.err); errCause != nil {
			jsonBody.ErrorCause = errCause
		}
	}

	if options.metadata != nil && len(options.metadata) > 0 {
		jsonBody.Metadata = options.metadata
	}

	if reqId, ok := getContextValue(options.context, RequestIdKey); ok {
		jsonBody.RequestId = reqId
	}

	// Compose json string
	jsonStr, _ := json.Marshal(jsonBody)
	l.writer.Printf("%s\n", jsonStr)
}

func getContextValue(ctx context.Context, key string) (string, bool) {
	if ctx == nil {
		return "", false
	}

	v, ok := ctx.Value(key).(string)
	return v, ok
}

type JsonBody struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	LevelId    int8                   `json:"levelId"`
	Message    string                 `json:"message"`
	Namespace  string                 `json:"namespace,omitempty"`
	Error      error                  `json:"error,omitempty"`
	ErrorCause error                  `json:"errorCause,omitempty"`
	Trace      string                 `json:"trace,omitempty"`
	RequestId  string                 `json:"requestId,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

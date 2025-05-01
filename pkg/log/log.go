package log

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

type Fields = logrus.Fields

func NewLogger() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetLevel(logrus.DebugLevel)

		logger.SetFormatter(&formatter.Formatter{
			NoColors:        false,
			TimestampFormat: "02 Jan 06 - 15:04",
			HideKeys:        false,
			CallerFirst:     true,
			CustomCallerFormatter: func(f *runtime.Frame) string {
				s := strings.Split(f.Function, ".")
				funcName := s[len(s)-1]
				return fmt.Sprintf(" \x1b[%dm[%s:%d][%s()]", 34, path.Base(f.File), f.Line, funcName)
			},
		})

		writers := []io.Writer{os.Stderr}

		appEnv := os.Getenv("APP_ENV")
		if appEnv != "test" {
			fileWriter := &lumberjack.Logger{
				Filename:   fmt.Sprintf("./storage/logs/app-%s.log", time.Now().Format("2006-01-02")),
				LocalTime:  true,
				Compress:   true,
				MaxSize:    100, // megabytes
				MaxAge:     7,   // days
				MaxBackups: 3,
			}
			writers = append(writers, fileWriter)
		}

		logger.SetOutput(io.MultiWriter(writers...))
		logger.SetReportCaller(true)
	})

	return logger
}

// Debug logs a message at the debug level with additional fields
func Debug(fields Fields, msg string) {
	if fields == nil {
		fields = Fields{}
	}
	logger.WithFields(fields).Debug(msg)
}

// Info logs a message at the info level with additional fields
func Info(fields Fields, msg string) {
	if fields == nil {
		fields = Fields{}
	}
	logger.WithFields(fields).Info(msg)
}

// Warn logs a message at the warning level with additional fields
func Warn(fields Fields, msg string) {
	if fields == nil {
		fields = Fields{}
	}
	logger.WithFields(fields).Warn(msg)
}

// Error logs a message at the error level with additional fields
func Error(fields Fields, msg string) {
	if fields == nil {
		fields = Fields{}
	}
	logger.WithFields(fields).Error(msg)
}

// ErrorWithTraceID logs an error and returns a trace ID for tracking purposes
func ErrorWithTraceID(fields Fields, msg string) uuid.UUID {
	traceID, err := uuid.NewRandom()
	if err != nil {
		Error(Fields{
			"error": err.Error(),
		}, "[log.ErrorWithTraceID] failed to generate trace ID")
	}

	if fields == nil {
		fields = Fields{}
	}
	fields["trace_id"] = traceID
	logger.WithFields(fields).Error(msg)

	return traceID
}

// Fatal logs a message at the fatal level with additional fields and terminates the program
func Fatal(fields Fields, msg string) {
	if fields == nil {
		fields = Fields{}
	}
	logger.WithFields(fields).Fatal(msg)
}

// Panic logs a message at the panic level with additional fields and panics
func Panic(fields Fields, msg string) {
	if fields == nil {
		fields = Fields{}
	}
	logger.WithFields(fields).Panic(msg)
}

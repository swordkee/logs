package logs

import (
	"f.in/v/logs/hooks/file"
	"fmt"
	"github.com/sirupsen/logrus"
	rsyslog "github.com/sirupsen/logrus/hooks/syslog"
	"io"
	"log/syslog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"f.in/v/logs/hooks/aliyun"
)

// Level describes the log severity level.
type Level uint8

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

// Logger is an interface that describes logging.
type Logger interface {
	SetLevel(level Level)
	SetOut(out io.Writer)

	Debug(...interface{})
	Debugln(...interface{})

	Info(...interface{})
	Infoln(...interface{})

	Warn(...interface{})
	Warnln(...interface{})

	Error(...interface{})
	Errorln(...interface{})

	Fatal(...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicln(...interface{})

	With(key string, value interface{}) Logger
	WithError(err error) Logger
}

type logger struct {
	entry *logrus.Entry
}

// With attaches a key-value pair to a logger.
func (l logger) With(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

// WithError attaches an error to a logger.
func (l logger) WithError(err error) Logger {
	return logger{l.entry.WithError(err)}
}

// SetLevel sets the level of a logger.
func (l logger) SetLevel(level Level) {
	l.entry.Logger.Level = logrus.Level(level)
}

// SetOut sets the output destination for a logger.
func (l logger) SetOut(out io.Writer) {
	l.entry.Logger.Out = out
}

func (l logger) SetHook(hookType, topic, network, logStore, accessKey, accessKeySecret string) {
	if hookType == "aliyun" {
		hook, err := aliyun.NewHook(network,
			accessKey, accessKeySecret, logStore, topic, false)
		if err != nil {
			l.entry.Logger.Error("Unable to connect to local aliyun daemon")
		} else {
			l.entry.Logger.AddHook(hook)
		}
	}
	if hookType == "syslog" {
		hook, err := rsyslog.NewSyslogHook("tcp", network, syslog.LOG_INFO, topic)
		if err != nil {
			l.entry.Logger.Error("Unable to connect to local syslog daemon")
		} else {
			l.entry.Logger.AddHook(hook)
		}
	} else if hookType == "files" {
		l.entry.Logger.AddHook(file.NewFileHook(selfDir() + "/log/" + topic + ".log"))
	}
}
func (l logger) SetFormat(format string) {
	if format == "text" {
		l.entry.Logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	} else {
		l.entry.Logger.Formatter = &logrus.JSONFormatter{}
	}
}

// Debug logs a message at level Debug on the standard logger.
func (l logger) Debug(args ...interface{}) {
	l.sourced().Debug(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func (l logger) Debugln(args ...interface{}) {
	l.sourced().Debugln(args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Info(args ...interface{}) {
	l.sourced().Info(args...)
}

// Infoln logs a message at level Info on the standard logger.
func (l logger) Infoln(args ...interface{}) {
	l.sourced().Infoln(args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warn(args ...interface{}) {
	l.sourced().Warn(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func (l logger) Warnln(args ...interface{}) {
	l.sourced().Warnln(args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Error(args ...interface{}) {
	l.sourced().Error(args...)
}

// Errorln logs a message at level Error on the standard logger.
func (l logger) Errorln(args ...interface{}) {
	l.sourced().Errorln(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatal(args ...interface{}) {
	l.sourced().Fatal(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func (l logger) Fatalln(args ...interface{}) {
	l.sourced().Fatalln(args...)
}

// Panic logs a message at level Panic on the standard logger.
func (l logger) Panic(args ...interface{}) {
	l.sourced().Panic(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func (l logger) Panicln(args ...interface{}) {
	l.sourced().Panicln(args...)
}

// sourced adds a source field to the logger that contains
// the file name and line where the logging happened.
func (l logger) sourced() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	fn := "(unknown)"
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		fn = runtime.FuncForPC(pc).Name()
	}
	logger := l.entry.WithField("source", fmt.Sprintf("%s:%d", file, line))
	return logger.WithField("source_func", fn)
}

var origLogger = logrus.New()
var baseLogger = logger{entry: logrus.NewEntry(origLogger)}

// New returns a new logger.
func New() Logger {
	return logger{entry: logrus.NewEntry(origLogger)}
}

// Base returns the base logger.
func Base() Logger {
	return baseLogger
}

// SetLevel sets the Level of the base logger
func SetLevel(level Level) {
	baseLogger.entry.Logger.Level = logrus.Level(level)
}

// SetOut sets the output destination base logger
func SetOut(out io.Writer) {
	baseLogger.entry.Logger.Out = out
}

func SetHook(hookType, topic, network, logStore, accessKey, accessKeySecret string) {
	if hookType == "aliyun" {
		hook, err := aliyun.NewHook(network, accessKey, accessKeySecret, logStore, topic, false)
		if err != nil {
			baseLogger.Error("Unable to connect to local aliyun daemon")
		} else {
			baseLogger.entry.Logger.AddHook(hook)
		}
	}
	if hookType == "syslog" {
		hook, err := rsyslog.NewSyslogHook("tcp", network, syslog.LOG_INFO, topic)
		if err != nil {
			baseLogger.Error("Unable to connect to local syslog daemon")
		} else {
			baseLogger.entry.Logger.AddHook(hook)
		}
	}
	if hookType == "files" {
		baseLogger.entry.Logger.AddHook(file.NewFileHook(selfDir() + "/log/" + topic + ".log"))
	}
}

func SetFormat(format string) {
	if format == "text" {
		baseLogger.entry.Logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	} else {
		baseLogger.entry.Logger.Formatter = &logrus.JSONFormatter{}
	}
}

// With attaches a key,value pair to a logger.
func With(key string, value interface{}) Logger {
	return baseLogger.With(key, value)
}

// WithError returns a Logger that will print an error along with the next message.
func WithError(err error) Logger {
	return logger{entry: baseLogger.sourced().WithError(err)}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	baseLogger.sourced().Debug(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	baseLogger.sourced().Debugln(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	baseLogger.sourced().Info(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	baseLogger.sourced().Infoln(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	baseLogger.sourced().Warn(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	baseLogger.sourced().Warnln(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	baseLogger.sourced().Error(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	baseLogger.sourced().Errorln(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	baseLogger.sourced().Fatal(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	baseLogger.sourced().Fatalln(args...)
}

// Panic logs a message at level Fatal on the standard logger.
func Panic(args ...interface{}) {
	baseLogger.sourced().Panic(args...)
}

// Panicln logs a message at level Fatal on the standard logger.
func Panicln(args ...interface{}) {
	baseLogger.sourced().Panicln(args...)
}

func selfDir() string {
	path, _ := filepath.Abs(os.Args[0])
	return filepath.Dir(path)
}

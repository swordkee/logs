package logs

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"runtime"
	"strings"
	"time"
)

type Logger interface {
	Debug(...interface{})
	Debugf(format string, args ...interface{})
	Debugln(...interface{})

	Info(...interface{})
	Infof(format string, args ...interface{})
	Infoln(...interface{})

	Warn(...interface{})
	Warnf(format string, args ...interface{})
	Warnln(...interface{})

	Error(...interface{})
	Errorf(format string, args ...interface{})
	Errorln(...interface{})

	Panic(...interface{})
	Panicf(format string, args ...interface{})
	Panicln(...interface{})

	Fatal(...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(...interface{})

	With(key string, value interface{}) *ZapLogger
}

// callerEncoder will add caller to log. format is "filename:lineNum:funcName", e.g:"logs/logs_test.go:15:logs.TestNew"
func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(strings.Join([]string{caller.TrimmedPath(), runtime.FuncForPC(caller.PC).Name()}, ":"))
}

// timeEncoder specifics the time format
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// milliSecondsDurationEncoder serializes a time.Duration to a floating-point number of micro seconds elapsed.
func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

func newLoggerConfig(debugLevel bool, te zapcore.TimeEncoder, de zapcore.DurationEncoder) (loggerConfig zap.Config) {
	loggerConfig = zap.NewProductionConfig()
	if te == nil {
		loggerConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		loggerConfig.EncoderConfig.EncodeTime = te
	}
	if de == nil {
		loggerConfig.EncoderConfig.EncodeDuration = milliSecondsDurationEncoder
	} else {
		loggerConfig.EncoderConfig.EncodeDuration = de
	}
	loggerConfig.EncoderConfig.EncodeCaller = callerEncoder
	if debugLevel {
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	return
}

// NewLogger return a normal logger
func NewLogger(debugLevel bool) (logger *zap.Logger) {
	loggerConfig := newLoggerConfig(debugLevel, nil, nil)
	loggerConfig.DisableStacktrace = true
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	return
}

// NewNoCallerLogger return a no caller key value, will be faster
func NewNoCallerLogger(debugLevel bool) (noCallerLogger *zap.Logger) {
	loggerConfig := newLoggerConfig(debugLevel, nil, nil)
	loggerConfig.DisableCaller = true
	noCallerLogger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	return
}

// CompatibleLogger is a logger which compatible to logrus/std log/prometheus.
// it implements Print() Println() Printf() Dbug() Debugln() Debugf() Info() Infoln() Infof() Warn() Warnln() Warnf()
// Error() Errorln() Errorf() Fatal() Fataln() Fatalf() Panic() Panicln() Panicf() With() WithField() WithFields()

type ZapLogger struct {
	_log *zap.Logger
}

// NewCompatibleLogger return CompatibleLogger with caller field
func NewCompatibleLogger(debugLevel bool) *ZapLogger {
	return &ZapLogger{NewLogger(debugLevel).WithOptions(zap.AddCallerSkip(1))}
}

// Debug logs a message at level Debug on the compatibleLogger.
func (l ZapLogger) Debug(args ...interface{}) {
	l._log.Debug(fmt.Sprint(args...))
}

// Debugln logs a message at level Debug on the compatibleLogger.
func (l ZapLogger) Debugln(args ...interface{}) {
	l._log.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at level Debug on the compatibleLogger.
func (l ZapLogger) Debugf(format string, args ...interface{}) {
	l._log.Debug(fmt.Sprintf(format, args...))
}

// Info logs a message at level Info on the compatibleLogger.
func (l ZapLogger) Info(args ...interface{}) {
	l._log.Info(fmt.Sprint(args...))
}

// Infoln logs a message at level Info on the compatibleLogger.
func (l ZapLogger) Infoln(args ...interface{}) {
	l._log.Info(fmt.Sprint(args...))
}

// Infof logs a message at level Info on the compatibleLogger.
func (l ZapLogger) Infof(format string, args ...interface{}) {
	l._log.Info(fmt.Sprintf(format, args...))
}

// Warn logs a message at level Warn on the compatibleLogger.
func (l ZapLogger) Warn(args ...interface{}) {
	l._log.Warn(fmt.Sprint(args...))
}

// Warnln logs a message at level Warn on the compatibleLogger.
func (l ZapLogger) Warnln(args ...interface{}) {
	l._log.Warn(fmt.Sprint(args...))
}

// Warnf logs a message at level Warn on the compatibleLogger.
func (l ZapLogger) Warnf(format string, args ...interface{}) {
	l._log.Warn(fmt.Sprintf(format, args...))
}

// Error logs a message at level Error on the compatibleLogger.
func (l ZapLogger) Error(args ...interface{}) {
	l._log.Error(fmt.Sprint(args...))
}

// Errorln logs a message at level Error on the compatibleLogger.
func (l ZapLogger) Errorln(args ...interface{}) {
	l._log.Error(fmt.Sprint(args...))
}

// Errorf logs a message at level Error on the compatibleLogger.
func (l ZapLogger) Errorf(format string, args ...interface{}) {
	l._log.Error(fmt.Sprintf(format, args...))
}

// Fatal logs a message at level Fatal on the compatibleLogger.
func (l ZapLogger) Fatal(args ...interface{}) {
	l._log.Fatal(fmt.Sprint(args...))
}

// Fatalln logs a message at level Fatal on the compatibleLogger.
func (l ZapLogger) Fatalln(args ...interface{}) {
	l._log.Fatal(fmt.Sprint(args...))
}

// Fatalf logs a message at level Fatal on the compatibleLogger.
func (l ZapLogger) Fatalf(format string, args ...interface{}) {
	l._log.Fatal(fmt.Sprintf(format, args...))
}

// Panic logs a message at level Painc on the compatibleLogger.
func (l ZapLogger) Panic(args ...interface{}) {
	l._log.Panic(fmt.Sprint(args...))
}

// Panicln logs a message at level Painc on the compatibleLogger.
func (l ZapLogger) Panicln(args ...interface{}) {
	l._log.Panic(fmt.Sprint(args...))
}

// Panicf logs a message at level Painc on the compatibleLogger.
func (l ZapLogger) Panicf(format string, args ...interface{}) {
	l._log.Panic(fmt.Sprintf(format, args...))
}

// With return a logger with an extra field.
func (l *ZapLogger) With(key string, value interface{}) *ZapLogger {
	return &ZapLogger{l._log.With(zap.Any(key, value))}
}

// WithField return a logger with an extra field.
func (l *ZapLogger) WithField(key string, value interface{}) *ZapLogger {
	return &ZapLogger{l._log.With(zap.Any(key, value))}
}

// WithFields return a logger with extra fields.
func (l *ZapLogger) WithFields(fields map[string]interface{}) *ZapLogger {
	i := 0
	var clog *ZapLogger
	for k, v := range fields {
		if i == 0 {
			clog = l.WithField(k, v)
		} else {
			clog = clog.WithField(k, v)
		}
		i++
	}
	return clog
}

// FormatStdLog set the output of stand package log to zaplog
func FormatStdLog() {
	log.SetFlags(log.Llongfile)
	log.SetOutput(&logWriter{NewNoCallerLogger(false)})
}

type logWriter struct {
	logger *zap.Logger
}

// Write implement io.Writer, as std log's output
func (w logWriter) Write(p []byte) (n int, err error) {
	i := bytes.Index(p, []byte(":")) + 1
	j := bytes.Index(p[i:], []byte(":")) + 1 + i
	caller := bytes.TrimRight(p[:j], ":")
	// find last index of /
	i = bytes.LastIndex(caller, []byte("/"))
	// find penultimate index of /
	i = bytes.LastIndex(caller[:i], []byte("/"))
	w.logger.Info("stdLog", zap.ByteString("caller", caller[i+1:]), zap.ByteString("log", bytes.TrimSpace(p[j:])))
	n = len(p)
	err = nil
	return
}

var baseLogger = &ZapLogger{NewLogger(false).WithOptions(zap.AddCallerSkip(2))}

// With attaches a key,value pair to a logger.
func With(key string, value interface{}) *ZapLogger {
	return baseLogger.With(key, value)
}

//// SetLevel sets the Level of the base logger
//func SetLevel(level zapcore.Level) {
//	baseLogger.SetLevel(level)
//}
//

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	baseLogger.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	baseLogger.Debug(fmt.Sprintf(format, args...))
}

// Debugw logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	baseLogger.Debug(fmt.Sprint(args...))
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	baseLogger.Info(fmt.Sprint(args...))
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	baseLogger.Info(fmt.Sprintf(format, args...))
}

// Infow logs a message at level Info on the standard logger.
func Infoln(format string, args ...interface{}) {
	baseLogger.Info(fmt.Sprintf(format, args...))
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	baseLogger.Warn(fmt.Sprint(args...))
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	baseLogger.Warnf(fmt.Sprintf(format, args...))
}

// Warnw logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	baseLogger.Warnln(fmt.Sprint(args...))
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	baseLogger.Error(fmt.Sprint(args...))
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	baseLogger.Error(fmt.Sprintf(format, args...))
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	baseLogger.Error(fmt.Sprint(args...))
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	baseLogger.Fatal(fmt.Sprint(args...))
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	baseLogger.Fatal(fmt.Sprintf(format, args...))
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	baseLogger.Fatal(fmt.Sprint(args...))
}

// Panic logs a message at level Fatal on the standard logger.
func Panic(args ...interface{}) {
	baseLogger.Panic(fmt.Sprint(args...))
}

// Panicf logs a message at level Fatal on the standard logger.
func Panicf(format string, args ...interface{}) {
	baseLogger.Panicf(fmt.Sprintf(format, args...))
}

// Panicln logs a message at level Fatal on the standard logger.
func Panicln(args ...interface{}) {
	baseLogger.Panic(fmt.Sprint(args...))
}

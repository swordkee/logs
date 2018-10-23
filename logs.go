package logs

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is an interface that describes logging.
type Logger interface {
	SetLevel(level zapcore.Level)

	Debug(...interface{})
	Debugf(format string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Info(...interface{})
	Infof(format string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Warn(...interface{})
	Warnf(format string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Error(...interface{})
	Errorf(format string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	DPanic(...interface{})
	DPanicf(format string, args ...interface{})

	Panic(...interface{})
	Panicf(format string, args ...interface{})

	Fatal(...interface{})
	Fatalf(format string, args ...interface{})

	With(...interface{}) *ZapLogger
}
type ProcessorFunc func() zapcore.Field

type ZapLogger struct {
	logger     *zap.Logger
	Config     zap.Config
	Level      zap.AtomicLevel
	Processors []ProcessorFunc
}

func NewLogger() (*ZapLogger, error) {
	a := zap.NewAtomicLevelAt(zap.InfoLevel)
	l := &ZapLogger{
		Config: zap.Config{
			Development:       true,
			DisableStacktrace: true,
			EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
			Encoding:          "json",
			ErrorOutputPaths:  []string{"stderr", "error.log"},
			Level:             a,
			OutputPaths:       []string{"stdout"},
		},
		Level: a,
	}

	logger, err := l.Config.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build logger")
	}
	defer logger.Sync()
	l.logger = logger
	return l, nil
}

// With attaches a fields pair to a logger.
//func (l *ZapLogger) With(args ...interface{}) *zap.SugaredLogger {
//	return l.logger.Sugar().With(args...)
//}
func (l *ZapLogger) With(fields ...zapcore.Field) *ZapLogger {
	n := &ZapLogger{
		logger:     l.logger.With(fields...),
		Config:     l.Config,
		Level:      l.Level,
		Processors: make([]ProcessorFunc, len(l.Processors)),
	}

	for i, p := range l.Processors {
		n.Processors[i] = p
	}
	return n
}

// SetLevel sets the level of a logger.
func (l *ZapLogger) SetLevel(level zapcore.Level) {
	l.Level.SetLevel(level)
}

//// SetOut sets the output destination for a logger.
//func (l *ZapLogger) SetOut(out io.Writer) {
//}
//

func (l *ZapLogger) SetFormat(format string) {
	l.Config.Encoding = format
	l.Config.Build()
}

// Debug logs a message at level Debug on the standard logger.
func (l *ZapLogger) Debug(args ...interface{}) {
	l.logger.Sugar().Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Sugar().Debugf(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (l *ZapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Debugw(msg, keysAndValues...)
}

// Info logs a message at level Info on the standard logger.
func (l *ZapLogger) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

// Infof logs a message at level Info on the standard logger.
func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.logger.Sugar().Infof(format, args...)
}

// Infow logs a message at level Info on the standard logger.
func (l *ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Infow(msg, keysAndValues...)
}

// Warn logs a message at level Warn on the standard logger.
func (l *ZapLogger) Warn(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Sugar().Warnf(format, args...)
}

// Warnw logs a message at level Warn on the standard logger.
func (l *ZapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Warnw(msg, keysAndValues...)
}

// Error logs a message at level Error on the standard logger.
func (l *ZapLogger) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Sugar().Errorf(format, args...)
}

// Errorw logs a message at level Error on the standard logger.
func (l *ZapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Errorw(msg, keysAndValues...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l *ZapLogger) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Sugar().Fatalf(format, args...)
}

// DPanic logs a message at level Panic on the standard logger.
func (l *ZapLogger) DPanic(args ...interface{}) {
	l.logger.Sugar().DPanic(args...)
}

// DPanicf logs a message at level Panic on the standard logger.
func (l *ZapLogger) DPanicf(format string, args ...interface{}) {
	l.logger.Sugar().DPanicf(format, args...)
}

// Panic logs a message at level Panic on the standard logger.
func (l *ZapLogger) Panic(args ...interface{}) {
	l.logger.Sugar().Panic(args...)
}

// Panicf logs a message at level Panic on the standard logger.
func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.logger.Sugar().Panicf(format, args...)
}

var baseLogger, _ = NewLogger()

// With attaches a key,value pair to a logger.
func With(fields ...zapcore.Field) *ZapLogger {
	return baseLogger.With(fields...)
}

// SetLevel sets the Level of the base logger
func SetLevel(level zapcore.Level) {
	baseLogger.SetLevel(level)
}

// SetFormat sets the Format of the base logger
func SetFormat(format string) {
	baseLogger.SetFormat(format)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	baseLogger.Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	baseLogger.Debugf(format, args...)
}

// Debugw logs a message at level Debug on the standard logger.
func Debugw(msg string, args ...interface{}) {
	baseLogger.Debugw(msg, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	baseLogger.Info(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	baseLogger.Infof(format, args...)
}

// Infow logs a message at level Info on the standard logger.
func Infow(msg string, keysAndValues ...interface{}) {
	baseLogger.Infow(msg, keysAndValues...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	baseLogger.Warn(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	baseLogger.Warnf(format, args...)
}

// Warnw logs a message at level Warn on the standard logger.
func Warnw(msg string, keysAndValues ...interface{}) {
	baseLogger.Warnw(msg, keysAndValues...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	baseLogger.Error(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	baseLogger.Errorf(format, args...)
}

// Errorw logs a message at level Error on the standard logger.
func Errorw(msg string, keysAndValues ...interface{}) {
	baseLogger.Errorw(msg, keysAndValues...)
}
func DPanic(args ...interface{}) {
	baseLogger.DPanic(args...)
}

// Panicln logs a message at level Fatal on the standard logger.
func DPanicf(format string, args ...interface{}) {
	baseLogger.DPanicf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	baseLogger.Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	baseLogger.Fatalf(format, args...)
}

// Panic logs a message at level Fatal on the standard logger.
func Panic(args ...interface{}) {
	baseLogger.Panic(args...)
}

// Panicln logs a message at level Fatal on the standard logger.
func Panicf(format string, args ...interface{}) {
	baseLogger.Panicf(format, args...)
}

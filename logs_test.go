package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"testing"
)

type user struct {
	Name string
}

func TestFormatStdLog(t *testing.T) {
	FormatStdLog()
	log.Print("std log")
	log.Print("")
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger(true)
	defer logger.Sync()
	logger.Debug("zap log debug msg")
}

func TestNewCompatibleLogger(t *testing.T) {
	compatibleLogger := NewCompatibleLogger(true)

	compatibleLogger.WithField("field", "value").Info("compatibleLogger with info")
	withFieldLogger := compatibleLogger.WithFields(map[string]interface{}{"field1": "value1", "field2": "value2"})
	withFieldLogger.Info("withFieldLogger Info")
	withFieldLogger.With("field3", "value3").Info("with filed3")
	withFieldLogger.Debugf("withFieldLogger debugf:%v", 1)
}

func BenchmarkStdLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		log.Print("std log printf")
	}
}

func BenchmarkNewLogger(b *testing.B) {
	logger := NewLogger(true)
	defer logger.Sync()
	for i := 0; i < b.N; i++ {
		logger.Debug("zap log debug msg")
	}
}
func TestNew(t *testing.T) {
	l := NewCompatibleLogger(false)
	//l.SetFormat("console")
	l.Info("hello world", zap.String("hello", "World"))
	l.Warn("info logging enabled")
	//l.SetLevel(zap.InfoLevel)
	l.Info("info logging disabled")
	l.Warnf("hello  %s, %d", "world", 12)
	l.With("ddd", "vvv").Info("dddd")
	//time.Sleep(time.Second * 5)

}
func Test(t *testing.T) {
	users := []*user{
		&user{Name: "Zap1"},
		&user{Name: "Zap2"},
		&user{Name: "Zap3"},
	}
	Warnf("hello  %s, %d", "world", 12)
	Debug("this is zap Debug")
	Debugf("this is zap %s %s", "test", "zap")
	Info("this is zap Info")
	////SetLevel(zap.ErrorLevel)
	Infof("hello  %s, %d", "world", 12)
	Errorln("hello111", "world", 12, "ttt", true)
	Warn("this is zap Warn")
	Warnf("hello %s,%s", "world", "red")
	Error("this is zap Error")

	//context := []interface{}{"foo", "bar"}
	//expectedFields := []zap.Field{zap.String("foo", "bar"), zap.Bool("baz", false)}

	Errorf("hello %s %s", "world", "red")

	Info("array sample", zap.Array("userArray", zapcore.ArrayMarshalerFunc(func(inner zapcore.ArrayEncoder) error {
		for _, u := range users {
			inner.AppendString(u.Name)
		}
		return nil
	})))
	Fatal("this is zap Fatal")
	//Fatalf("hello %s", "world", "red")
	Panic("this is zap Panic")
	Panicf("hello %s ,%s", "world", "red")
	//time.Sleep(time.Second * 5)

}

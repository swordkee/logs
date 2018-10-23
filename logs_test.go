package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

type user struct {
	Name string
}

func TestNew(t *testing.T) {
	l, _ := NewLogger()
	l.SetFormat("console")
	l.Info("hello world", zap.String("hello", "World"))
	l.Warn("info logging enabled")
	l.SetLevel(zap.InfoLevel)
	l.Info("info logging disabled")
	//time.Sleep(time.Second * 5)

}
func Test(t *testing.T) {
	users := []*user{
		&user{Name: "Zap1"},
		&user{Name: "Zap2"},
		&user{Name: "Zap3"},
	}
	Debug("this is zap Debug")
	Debugf("this is zap %s %s", "test", "zap")
	Info("this is zap Info")
	//SetLevel(zap.ErrorLevel)
	Infof("hello  %s, %d", "world", 12)
	Errorw("hello111", "world", 12, "ttt", true)
	Warn("this is zap Warn")
	Warnf("hello %s", "world", "red")
	Error("this is zap Error")

	//With(
	//	"hello111111111", "world",
	//	"failure", errors.New("oh no"),
	//	//Stack(),
	//	"count", 42,
	//	"user", user{Name: "alice"},
	//)

	context := []interface{}{"foo", "bar"}
	expectedFields := []zap.Field{zap.String("foo", "bar"), zap.Bool("baz", false)}
	With(expectedFields...).Info(context)

	Errorf("hello %s %s", "world", "red")

	Info("array sample", zap.Array("userArray", zapcore.ArrayMarshalerFunc(func(inner zapcore.ArrayEncoder) error {
		for _, u := range users {
			inner.AppendString(u.Name)
		}
		return nil
	})))
	Fatal("this is zap Fatal")
	Fatalf("hello %s", "world", "red")
	Panic("this is zap Panic")
	Panicf("hello %d", "world", "red")
	//time.Sleep(time.Second * 5)

}

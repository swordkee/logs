package file

import (
	"fmt"
	"os"
	"strings"
	"github.com/sirupsen/logrus"
)

func NewFileHook(file string) (f *FileHook) {
	w := NewLogFile(file)
	return &FileHook{w}
}

type FileHook struct {
	W LoggerInterface
}

func (hook *FileHook) Fire(entry *logrus.Entry) (err error) {
	message, err := getMessage(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	switch entry.Level {
	case logrus.PanicLevel:
		return hook.W.WriteMsg(message, LevelPanic)
	case logrus.FatalLevel:
		return hook.W.WriteMsg(message, LevelFatal)
	case logrus.ErrorLevel:
		return hook.W.WriteMsg(message, LevelError)
	case logrus.WarnLevel:
		return hook.W.WriteMsg(message, LevelWarn)
	case logrus.InfoLevel:
		return hook.W.WriteMsg(message, LevelInfo)
	case logrus.DebugLevel:
		return hook.W.WriteMsg(message, LevelDebug)
	default:
		return nil
	}
	return
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func getMessage(entry *logrus.Entry) (message string, err error) {
	message = message + fmt.Sprintf("%s", entry.Message)
	for k, v := range entry.Data {
		if !strings.HasPrefix(k, "err_") {
			message = message + fmt.Sprintf(" [%v] %v", k, v)
		}
	}
	if full, ok := entry.Data["err_full"]; ok {
		message = message + fmt.Sprintf("%v", full)
	} else {
		//file, lineNumber := caller.GetCallerIgnoringLogMulti(2)
		//if file != "" {
		//	sep := fmt.Sprintf("%s/src/", os.Getenv("GOPATH"))
		//	fileName := strings.Split(file, sep)
		//	if len(fileName) >= 2 {
		//		file = fileName[1]
		//	}
		//}
		//message = message + fmt.Sprintf("%s:%d", file, lineNumber)
	}

	return
}

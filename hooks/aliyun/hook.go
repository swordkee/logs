package aliyun

import (
	"github.com/sirupsen/logrus"
	"fmt"
	"os"
)

func NewHook(network, accessKey, accessSecret, logStore, topic string, isAsync bool) (*AliYunHook, error) {
	writer, err := NewWriter(network, accessKey, accessSecret, logStore, topic, isAsync)
	return &AliYunHook{writer}, err
}

type AliYunHook struct {
	w *Writer
}

func (hook *AliYunHook) Fire(entry *logrus.Entry) (err error) {
	logrus.SetOutput(hook.w)
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	_, err = hook.w.Write([]byte(line))
	return err
}

func (hook *AliYunHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

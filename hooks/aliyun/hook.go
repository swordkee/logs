package aliyun

import (
	"github.com/sirupsen/logrus"
	"context"
	"io"
	"fmt"
	"os"
)

type WriterMap map[logrus.Level]io.Writer

func NewHook(network, accessKey, accessSecret, logStore, topic string) (*AliYunHook, error) {
	ctx, _ := context.WithCancel(context.Background())
	writer, err := NewWriter(network, accessKey, accessSecret, logStore, topic, ctx)
	return &AliYunHook{writer}, err
}

type AliYunHook struct {
	w *Writer
}

func (hook *AliYunHook) Fire(entry *logrus.Entry) (err error) {
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

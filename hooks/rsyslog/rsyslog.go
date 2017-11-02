package rsyslog

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	"strings"
)

// SyslogHook to send logs via syslog.
type RSyslogHook struct {
	Writer *syslog.Writer
}

func NewSyslogHook(network, raddr string, priority syslog.Priority, tag string) (*RSyslogHook, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &RSyslogHook{w}, err
}

func (hook *RSyslogHook) Fire(entry *logrus.Entry) error {
	message, err := getMessage(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}
	switch entry.Level {
	case logrus.PanicLevel:
		return hook.Writer.Crit("[PANIC] " + message)
	case logrus.FatalLevel:
		return hook.Writer.Crit("[FATAL] " + message)
	case logrus.ErrorLevel:
		return hook.Writer.Err("[ERROR] " + message)
	case logrus.WarnLevel:
		return hook.Writer.Warning("[WARN] " + message)
	case logrus.InfoLevel:
		return hook.Writer.Info("[INFO] " + message)
	case logrus.DebugLevel:
		return hook.Writer.Debug("[DEBUG] " + message)
	default:
		return nil
	}
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
	}

	return
}

func (hook *RSyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

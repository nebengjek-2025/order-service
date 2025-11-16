package log

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Log struct singleton
type Log struct {
	AppName  string
	LogLevel int
	Logger   *logrus.Logger
	// future: add logstash/kafka client fields here
}

var logger Log

var mapOfLogLevel = map[string]int{
	"DEBUG": 1,
	"ERROR": 2,
}

// InitLogger initialize logger from Viper
func InitLogger(v *viper.Viper) {
	levelStr := v.GetString("log.level")
	appName := v.GetString("app.name")

	logger = Log{
		AppName:  appName,
		LogLevel: mapOfLogLevel[levelStr],
		Logger:   newLogrusLogger(v),
	}
}

// GetLogger return singleton
func GetLogger() Log {
	return logger
}

// internal helper to create logrus instance
func newLogrusLogger(v *viper.Viper) *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	levelStr := v.GetString("log.level")
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.InfoLevel
	}
	l.SetLevel(level)
	return l
}

// -----------------------------
// Info
func (l Log) Info(context, message, scope, meta string) {
	if l.LogLevel <= 1 {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("[INFO] Service: %s - Context: %s - Message: %s - Scope: %s - Meta: %s - At: %s:%d",
			l.AppName, context, message, scope, meta, file, line)
		println(msg)
		l.Logger.WithFields(logrus.Fields{
			"service": l.AppName,
			"context": context,
			"scope":   scope,
			"meta":    meta,
			"file":    file,
			"line":    line,
		}).Info(message)
	}
}

// -----------------------------
// Error
func (l Log) Error(context, message, scope, meta string) {
	if l.LogLevel <= 2 {
		_, file, line, _ := runtime.Caller(1)
		_, file2, line2, _ := runtime.Caller(2)
		msg := fmt.Sprintf("[ERROR] Context: %s - Message: %s - Scope: %s - Meta: %s - At Level1: %s:%d - At Level2: %s:%d",
			context, message, scope, meta, file, line, file2, line2)
		println(msg)
		l.Logger.WithFields(logrus.Fields{
			"service": l.AppName,
			"context": context,
			"scope":   scope,
			"meta":    meta,
			"file1":   file,
			"line1":   line,
			"file2":   file2,
			"line2":   line2,
		}).Error(message)
	}
}

// -----------------------------
// Slow
func (l Log) Slow(context, message, scope, meta string) {
	if l.LogLevel <= 1 {
		_, file, line, _ := runtime.Caller(2)
		msg := fmt.Sprintf("[SLOW] Context: %s - Message: %s - Scope: %s - Meta: %s - At: %s:%d",
			context, message, scope, meta, file, line)
		println(msg)
		l.Logger.WithFields(logrus.Fields{
			"service": l.AppName,
			"context": context,
			"scope":   scope,
			"meta":    meta,
			"file":    file,
			"line":    line,
		}).Info("[SLOW] " + message)
	}
}

// -----------------------------
// Optional: send to logstash/kafka
func (l Log) SendTo3rdParty(level, msg, context, scope, meta, file string, line int) {
	payload := map[string]interface{}{
		"service": l.AppName,
		"level":   level,
		"msg":     msg,
		"context": context,
		"scope":   scope,
		"meta":    meta,
		"file":    file,
		"line":    line,
		"time":    time.Now().Format(time.RFC3339),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		l.Logger.WithError(err).Error("Failed marshal log payload")
		return
	}
	// TODO: kirim ke Logstash/Kafka
	_ = data
}

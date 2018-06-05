package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/FactomProject/logrustash"
	"time"
	"strings"
	"os"
)

const Production = "production"

type logger struct {
	logger *logrus.Logger
}

type entry struct {
	entry *logrus.Entry
}

type Logger interface {
	WithFields(map[string]interface{}) Logger
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Fatal(string)
	Panic(string)
}

type Config struct {
	ServiceName     string
	AppName         string
	EnvironmentName string
	LogstashServer  string
	LogstashPort    string
	LogLevel        string
}

func NewLogger(config *Config) Logger {
	newLogger := &logger{
		logger: logrus.StandardLogger(),
	}
	configure(config)
	return newLogger

}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &entry{
		entry: l.logger.WithFields(fields),
	}
}

func (l *logger) Debug(message string) {
	l.logger.Debug(message)
}

func (l *logger) Info(message string) {
	l.logger.Info(message)
}

func (l *logger) Warn(message string) {
	l.logger.Warn(message)
}

func (l *logger) Error(message string) {
	l.logger.Error(message)
}

func (l *logger) Fatal(message string) {
	l.logger.Fatal(message)
}

func (l *logger) Panic(message string) {
	l.logger.Panic(message)
}

func (e *entry) WithFields(fields map[string]interface{}) Logger {
	e.entry.WithFields(fields)
	return e
}

func (e *entry) Debug(message string) {
	e.entry.Debug(message)
}

func (e *entry) Info(message string) {
	e.entry.Info(message)
}

func (e *entry) Warn(message string) {
	e.entry.Warn(message)
}

func (e *entry) Error(message string) {
	e.entry.Error(message)
}

func (e *entry) Fatal(message string) {
	e.entry.Fatal(message)
}

func (e *entry) Panic(message string) {
	e.entry.Panic(message)
}

func configure(configuration *Config) {

	logrus.SetLevel(getLevel(configuration.LogLevel))
	logrus.SetFormatter(getFormatter(configuration.EnvironmentName))

	if hook := createHook(configuration); hook != nil {
		logrus.AddHook(hook)
	}
}

func getLevel(logLevel string) logrus.Level {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return logrus.InfoLevel
	}
	return level
}

func getFormatter(environment string) logrus.Formatter {
	envType := valueOrDefault(environment, "development")

	if envType == Production {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{}
}

func createHook(configuration *Config) logrus.Hook {
	logstashServer := configuration.LogstashServer
	if strings.TrimSpace(logstashServer) == "" {
		return nil
	}

	logstashPort := valueOrDefault(configuration.LogstashPort, "5000")
	appName := valueOrDefault(configuration.AppName, configuration.ServiceName)
	address := logstashServer + ":" + logstashPort

	hook, err := logrustash.NewAsyncHook("tcp", address, appName)

	if err != nil {
		return nil
	}

	hook.ReconnectBaseDelay = time.Second
	hook.ReconnectDelayMultiplier = 2
	hook.MaxReconnectRetries = 10

	return hook
}

func valueOrDefault(name string, defValue string) string {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue
	} else {
		return v
	}
}

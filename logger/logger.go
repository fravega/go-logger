package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

const json = "JSON"

//Alias for map used in withfields methods
type Fields map[string]interface{}

type logger struct {
	logger  *logrus.Logger
	dFields map[string]interface{}
}

var defaultLogger Logger
var once sync.Once

type entry struct {
	entry   *logrus.Entry
	dFields map[string]interface{}
}

type Logger interface {
	WithFields(map[string]interface{}) Logger
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Panicf(string, ...interface{})
}

type Config struct {
	ServiceName     string
	EnvironmentName string
	LogLevel        string
	Format          string
	DefaultFields   map[string]interface{}
}

func New(config *Config) Logger {
	fields := addIfNotEmpty(config.DefaultFields, "service_name", config.ServiceName)
	fields = addIfNotEmpty(fields, "environment", config.EnvironmentName)
	newLogger := &logger{
		logger:  logrus.StandardLogger(),
		dFields: fields,
	}
	configure(config)
	return newLogger

}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &entry{
		entry:   l.logger.WithFields(collectFields(l.dFields, fields)),
		dFields: l.dFields,
	}
}

func (l *logger) Debug(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Debug(message)
}

func (l *logger) Info(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Info(message)
}

func (l *logger) Warn(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Warn(message)
}

func (l *logger) Error(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Error(message)
}

func (l *logger) Fatal(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Fatal(message)
}

func (l *logger) Panic(message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Panic(message)
}

func (l *logger) Debugf(format string, message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Debugf(format, message...)
}

func (l *logger) Infof(format string, message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Infof(format, message...)
}

func (l *logger) Warnf(format string, message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Warnf(format, message...)
}

func (l *logger) Errorf(format string, message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Errorf(format, message...)
}

func (l *logger) Fatalf(format string, message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Fatalf(format, message...)
}

func (l *logger) Panicf(format string, message ...interface{}) {
	l.logger.WithFields(collectFields(l.dFields, map[string]interface{}{})).Panicf(format, message...)
}

func (e *entry) WithFields(fields map[string]interface{}) Logger {
	e.entry.WithFields(collectFields(e.dFields, fields))
	return e
}

func (e *entry) Debug(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Debug(message)
}

func (e *entry) Info(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Info(message)
}

func (e *entry) Warn(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Warn(message)
}

func (e *entry) Error(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Error(message)
}

func (e *entry) Fatal(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Fatal(message)
}

func (e *entry) Panic(message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Panic(message)
}

func (e *entry) Debugf(format string, message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Debugf(format, message...)
}

func (e *entry) Infof(format string, message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Infof(format, message...)
}

func (e *entry) Warnf(format string, message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Warnf(format, message...)
}

func (e *entry) Errorf(format string, message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Errorf(format, message...)
}

func (e *entry) Fatalf(format string, message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Fatalf(format, message...)
}

func (e *entry) Panicf(format string, message ...interface{}) {
	e.entry.WithFields(collectFields(e.dFields, map[string]interface{}{})).Panicf(format, message...)
}

func collectFields(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	var allFields = make(map[string]interface{}, len(a)+len(b))
	for k, v := range a {
		allFields[k] = v
	}
	for k, v := range b {
		allFields[k] = v
	}
	return allFields
}

func configure(configuration *Config) {
	logrus.SetLevel(getLevel(configuration.LogLevel))
	logrus.SetFormatter(getFormatter(configuration.Format))
}

func getLevel(logLevel string) logrus.Level {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return logrus.InfoLevel
	}
	return level
}

func getFormatter(format string) logrus.Formatter {
	envType := valueOrDefault(format, "plain")

	if strings.ToUpper(envType) == json {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{}
}

func valueOrDefault(name string, defValue string) string {
	v := strings.TrimSpace(os.Getenv(name))

	if v == "" {
		return defValue
	} else {
		return v
	}
}

func GetDefaultLogger() Logger {
	once.Do(func() {
		defaultLogger = buildDefaultLogger()
	})
	return defaultLogger
}

func buildDefaultLogger() Logger {
	config := &Config{
		ServiceName:     os.Getenv("SERVICE_NAME"),
		EnvironmentName: os.Getenv("ENVIRONMENT"),
		LogLevel:        os.Getenv("LOG_LEVEL"),
	}
	return New(config)
}

func addIfNotEmpty(fields map[string]interface{}, key string, value string) map[string]interface{} {
	if strings.TrimSpace(value) != "" {
		newFields := fields
		if len(newFields) == 0 {
			newFields = make(map[string]interface{})
		}
		newFields[key] = value
		return newFields
	} else {
		return fields
	}
}

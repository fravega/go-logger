package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/FactomProject/logrustash"
	"time"
	"strings"
	"os"
)

type Level uint32

const (
	Panic Level = iota
	Fatal
	Error
	Warn
	Info
	Debug
)

const Production = "production"

type logger struct {
	logger *logrus.Logger
}

type Logger interface {
	LogWithFields(Level, map[string]interface{}, string)
	Log(Level, string)
}

type Config struct {
	ServiceName     string
	AppName         string
	EnvironmentName string
	LogstashServer  string
	LogstashPort    string
	LogLevel        string
}

func NewLogger(config Config) Logger {
	newLogger := &logger{
		logger: logrus.StandardLogger(),
	}
	configure(&config)
	return newLogger

}

func (l *logger) LogWithFields(level Level, fields map[string]interface{}, message string){
	if strings.TrimSpace(message) == ""{
		return
	}

	withFields := l.logger.WithFields(fields)

	switch level {
	case Panic:
		withFields.Panic(message)
	case Fatal:
		withFields.Fatal(message)
	case Error:
		withFields.Error(message)
	case Warn:
		withFields.Warn(message)
	case Info:
		withFields.Info(message)
	case Debug:
		withFields.Debug(message)
	}
}

func (l *logger) Log(level Level, message string){
	if strings.TrimSpace(message) == ""{
		return
	}

	switch level {
	case Panic:
		l.logger.Panic(message)
	case Fatal:
		l.logger.Fatal(message)
	case Error:
		l.logger.Error(message)
	case Warn:
		l.logger.Warn(message)
	case Info:
		l.logger.Info(message)
	case Debug:
		l.logger.Debug(message)
	}
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
	if logstashServer == "" {
		return nil
	}

	logstashPort := valueOrDefault(configuration.LogstashPort, "5000")
	appName := valueOrDefault(configuration.AppName, configuration.ServiceName)
	address := logstashServer + ":" + logstashPort

	hook, err := logrustash.NewAsyncHook("tcp", address, appName)

	if err != nil {
		return nil
	}

	hook.ReconnectBaseDelay = time.Second // Wait for one second before first reconnect.
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

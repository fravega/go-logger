package logger

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Debug(t *testing.T) {
	testLoggerLevel(t, New(buildDefaultConfig()), "debug")
}

func TestLogger_Info(t *testing.T) {
	testLoggerLevel(t, New(buildDefaultConfig()), "info")
}

func TestLogger_Warn(t *testing.T) {
	testLoggerLevel(t, New(buildDefaultConfig()), "warn")
}

func TestLogger_Error(t *testing.T) {
	testLoggerLevel(t, New(buildDefaultConfig()), "error")
}

func TestLogger_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Log with level panic didn't panic")
		}
	}()

	testLoggerLevel(t, New(buildDefaultConfig()), "panic")
}

func testLoggerLevel(t *testing.T, sut Logger, level string) {

	var output bytes.Buffer
	logrus.SetOutput(&output)

	reflect.ValueOf(sut).MethodByName(strings.Title(level)).Call([]reflect.Value{reflect.ValueOf("sarasa")})

	assert.Contains(t, output.String(), fmt.Sprintf("level=%s", level), fmt.Sprintf("log entry is in %s level", level))
}

func TestLogger_LogWithDefaultFields(t *testing.T) {

	config := buildDefaultConfig()
	msg := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

	var output bytes.Buffer
	logrus.SetOutput(&output)

	sut := New(config)

	sut.Debugf(msg)

	assertForLogEntry(t, output.String(), msg, config.DefaultFields)
}

func TestLogger_LogWithAdditionalFields(t *testing.T) {

	config := buildDefaultConfig()
	msg := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

	var output bytes.Buffer
	logrus.SetOutput(&output)

	fields := map[string]interface{}{
		"customField3": true,
	}

	sut := New(config)
	sut.WithFields(fields).Infof(msg)

	assertForLogEntry(t, output.String(), msg, collectFields(config.DefaultFields, fields))
}
func TestLogger_LogWithAdditionalFieldsAndDifferentInstances(t *testing.T) {

	config := buildDefaultConfig()
	msg := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

	var output bytes.Buffer
	logrus.SetOutput(&output)

	fields := map[string]interface{}{
		"customField3": true,
	}

	sut := New(config)
	sut = sut.WithFields(fields)

	sut.Warnf(msg)

	assertForLogEntry(t, output.String(), msg, collectFields(config.DefaultFields, fields))

	output.Reset()

	fields2 := map[string]interface{}{
		"customField4": 4,
	}

	sut = sut.WithFields(fields2)

	sut.Errorf(msg)

	assertForLogEntry(t, output.String(), msg, collectFields(collectFields(config.DefaultFields, fields), fields2))
}

func assertForLogEntry(t *testing.T, logEntry string, expectedMsg string, expectedFields map[string]interface{}) {

	assert.Contains(t, logEntry, expectedMsg, "log entry contains specified message")
	for k, v := range expectedFields {
		assert.Contains(t, logEntry, fmt.Sprintf("%s=%v", k, v), fmt.Sprintf("log entry contains field %s", k))
	}
}

func buildDefaultConfig() *Config {
	return &Config{
		LogLevel:        "DEBUG",
		EnvironmentName: "TEST",
		ServiceName:     "LOG_TEST",
		DefaultFields: map[string]interface{}{
			"customField1": 1,
			"customField2": "2",
		},
	}
}

func TestLogger_Log(t *testing.T) {

	config := &Config{
		LogLevel:        "DEBUG",
		EnvironmentName: "TEST",
		ServiceName:     "LOG_TEST",
		DefaultFields: map[string]interface{}{
			"customField": 1,
		},
	}

	logger := New(config)

	var firstLog bytes.Buffer

	logrus.SetOutput(&firstLog)

	logger.Debug("Hola Logger")

	firstLogResult := firstLog.String()

	println("first: " + firstLogResult)

	if !strings.Contains(firstLogResult, "Hola Logger") {
		t.Fail()
	}

	if !strings.Contains(firstLogResult, "customField=1") {
		t.Fail()
	}

	fields := make(map[string]interface{})
	fields["text"] = "Text"
	fields["num"] = 2
	fields["obj"] = struct {
		name string
	}{
		name: "Name",
	}

	var secondLog bytes.Buffer

	logrus.SetOutput(&secondLog)

	logger.WithFields(fields).Debug("Hola Logger")

	secondLogResult := secondLog.String()

	println("second: " + secondLogResult)

	if !strings.Contains(secondLogResult, "num=2") {
		println("Fail on num")
		t.Fail()
	}

	if !strings.Contains(secondLogResult, "obj=\"{Name}\"") {
		println("Fail on obj")
		t.Fail()
	}

	if !strings.Contains(secondLogResult, "text=Text") {
		println("Fail on text")
		t.Fail()
	}

	emptyFields := make(map[string]interface{})

	var thirdLog bytes.Buffer

	logrus.SetOutput(&thirdLog)

	logger.WithFields(emptyFields).Debug("Hola Logger")

	thirdLogResult := thirdLog.String()

	println("third: " + thirdLogResult)

	if !strings.Contains(firstLogResult, "Hola Logger") {
		t.Fail()
	}

	if !strings.Contains(firstLogResult, "customField=1") {
		t.Fail()
	}
}

func TestLogger_GetDefaultLogger(t *testing.T) {

	os.Setenv("SERVICE_NAME", "katalog-service")
	os.Setenv("LOG_LEVEL", "INFO")

	logger := GetDefaultLogger()

	if logger != GetDefaultLogger() {
		t.Fail()
	}

	var firstLog bytes.Buffer

	logrus.SetOutput(&firstLog)

	logger.Info("First message")

	firstLogResult := firstLog.String()

	println("first: " + firstLogResult)

	if !strings.Contains(firstLogResult, "service_name=katalog-service") {
		t.Fail()
	}

	var secondLog bytes.Buffer

	logrus.SetOutput(&secondLog)

	logger.Debug("Second message")

	secondLogResult := secondLog.String()

	println("second: " + secondLogResult)

	if secondLogResult != "" {
		t.Fail()
	}

	fields := make(map[string]interface{})
	fields["mode"] = "application"

	logger.WithFields(fields).Warn("Second message")

	thirdLogResult := secondLog.String()

	println("third: " + thirdLogResult)

	if !strings.Contains(thirdLogResult, "mode=application") {
		t.Fail()
	}

}

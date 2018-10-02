package logger

import (
	"testing"
	"github.com/sirupsen/logrus"
	"bytes"
	"strings"
	"os"
)

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

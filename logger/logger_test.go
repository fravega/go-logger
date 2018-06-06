package logger

import (
	"testing"
	"github.com/sirupsen/logrus"
	"bytes"
	"strings"
)

func TestLogger_Log(t *testing.T) {

	config := &Config{
		LogLevel:        "DEBUG",
		EnvironmentName: "TEST",
		ServiceName:     "LOG_TEST",
		AppName:         "GOLANG-LOGGER",
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

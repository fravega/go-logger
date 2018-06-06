package logger

import (
	"testing"
	"github.com/sirupsen/logrus"
	"bytes"
	"strings"
)

func TestLogger_Log(t *testing.T) {

	config := &Config{
		LogLevel:"DEBUG",
		EnvironmentName:"TEST",
		ServiceName:"LOG_TEST",
		AppName:"GOLANG-LOGGER",
	}

	logger := New(config)

	var firstLog bytes.Buffer

	logrus.SetOutput(&firstLog)

	logger.Debug("Hola Logger")

	if !strings.Contains(firstLog.String(), "Hola Logger"){
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

	println("second: "+secondLogResult)

	if !strings.Contains(secondLogResult,"num=2"){
		println("Fail on num")
		t.Fail()
	}

	if !strings.Contains(secondLogResult,"obj=\"{Name}\""){
		println("Fail on obj")
		t.Fail()
	}

	if !strings.Contains(secondLogResult,"text=Text"){
		println("Fail on text")
		t.Fail()
	}

	emptyFields := make(map[string]interface{})

	logger.WithFields(emptyFields).Debug("Hola Logger")
}

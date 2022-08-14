package sterrors_logrus_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/Bibob7/sterrors"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

type TestWriter struct {
	Content []byte
}

func (w *TestWriter) Write(p []byte) (n int, err error) {
	w.Content = p
	return len(p), nil
}

type TestOutput struct {
	Msg       string               `json:"msg"`
	Level     string               `json:"level"`
	CallStack []TestCallStackEntry `json:"stack"`
}

type TestCallStackEntry struct {
	ErrMessage string          `json:"msg"`
	Severity   string          `json:"severity"`
	Caller     sterrors.Caller `json:"caller"`
}

func TestLogrusFormatter_Log(t *testing.T) {
	outputWriter := &TestWriter{}
	logrus.SetOutput(outputWriter)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	err := sterrors.E("initial error", sterrors.SeverityWarning)
	secondErr := sterrors.E("second error", err, sterrors.SeverityError)

	formatter := LogrusLogger{}
	formatter.Log(secondErr)

	var output TestOutput
	unmarshalErr := json.Unmarshal(outputWriter.Content, &output)

	e := err.(sterrors.Error)
	secondE := secondErr.(sterrors.Error)

	if unmarshalErr != nil {
		t.Errorf("expected nil, but got: %#v", unmarshalErr)
	}
	expectedMsg := "sterrors.TestLogrusFormatter_Log: second error"
	if expectedMsg != output.Msg {
		t.Errorf("expected message %s is not equal with actual message %s", expectedMsg, output.Msg)
	}
	if "error" != output.Level {
		t.Errorf("level ist not equal with error")
	}
	if len(output.CallStack) != 2 {
		t.Errorf("callstack length is not 2")
	}
	if !reflect.DeepEqual(TestCallStackEntry{
		ErrMessage: "second error",
		Severity:   "error",
		Caller: sterrors.Caller{
			FuncName: secondE.Caller().FuncName,
			File:     secondE.Caller().File,
			Line:     secondE.Caller().Line,
		},
	}, output.CallStack[0]) {
		t.Errorf("callstack entry is not equal with expected entry")
	}
	if !reflect.DeepEqual(TestCallStackEntry{
		ErrMessage: "initial error",
		Severity:   "warning",
		Caller: sterrors.Caller{
			FuncName: e.Caller().FuncName,
			File:     e.Caller().File,
			Line:     e.Caller().Line,
		},
	}, output.CallStack[1]) {
		t.Errorf("callstack entry is not equal with expected entry")
	}
}

func TestLogrusLogger_LogHighestSeverities(t *testing.T) {
	testCases := map[string]struct {
		HighestSeverity sterrors.Severity
		LogrusLevel     string
	}{
		"Severity Unknown": {
			HighestSeverity: -1,
			LogrusLevel:     "error",
		},
		"Severity Error": {
			HighestSeverity: sterrors.SeverityError,
			LogrusLevel:     "error",
		},
		"Severity Warning": {
			HighestSeverity: sterrors.SeverityWarning,
			LogrusLevel:     "warning",
		},
		"Severity Info": {
			HighestSeverity: sterrors.SeverityInfo,
			LogrusLevel:     "info",
		},
		"Severity Notice": {
			HighestSeverity: sterrors.SeverityNotice,
			LogrusLevel:     "info",
		},
		"Severity Debug": {
			HighestSeverity: sterrors.SeverityDebug,
			LogrusLevel:     "debug",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outputWriter := &TestWriter{}
			logrus.SetOutput(outputWriter)
			logrus.SetFormatter(&logrus.JSONFormatter{})
			logrus.SetLevel(logrus.DebugLevel)

			var err error
			if tc.HighestSeverity == -1 {
				err = fmt.Errorf("initial error")
			} else {
				err = sterrors.E("initial error", tc.HighestSeverity)
			}

			formatter := LogrusLogger{}
			formatter.Log(err)

			fmt.Println(string(outputWriter.Content))
			var output TestOutput
			unmarshalErr := json.Unmarshal(outputWriter.Content, &output)
			if unmarshalErr != nil {
				t.Errorf("expected nil, but got: %#v", unmarshalErr)
			}
			if !reflect.DeepEqual(tc.LogrusLevel, output.Level) {
				t.Errorf("expected level %v is not equal with %v", tc.LogrusLevel, output.Level)
			}
		})
	}
}

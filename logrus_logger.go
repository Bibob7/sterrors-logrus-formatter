package sterrors_logrus_plugin

import (
	"github.com/Bibob7/sterrors"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct{}

func (f *LogrusLogger) Log(err error) {
	e, ok := err.(sterrors.Error)
	if !ok {
		logrus.Error(err)
		return
	}
	entry := logrus.WithFields(logrus.Fields{"stack": sterrors.CallStack(e)})
	switch sterrors.HighestSeverity(e) {
	case sterrors.SeverityInfo:
		entry.Infof("%s: %v", e.Caller().FuncName, e)
	case sterrors.SeverityNotice:
		entry.Infof("%s: %v", e.Caller().FuncName, e)
	case sterrors.SeverityWarning:
		entry.Warnf("%s: %v", e.Caller().FuncName, e)
	case sterrors.SeverityDebug:
		entry.Debugf("%s: %v", e.Caller().FuncName, e)
	default:
		entry.Errorf("%s: %v", e.Caller().FuncName, e)
	}
}

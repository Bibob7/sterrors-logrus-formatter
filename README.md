# Sterrors

[![codecov](https://codecov.io/gh/Bibob7/sterrors-logrus-plugin/branch/main/graph/badge.svg?token=2LURD0VD9X)](https://codecov.io/gh/Bibob7/sterrors-logrus-plugin)
[![CircleCI](https://circleci.com/gh/Bibob7/sterrors-logrus-plugin/tree/main.svg?style=svg)](https://circleci.com/gh/Bibob7/sterrors-logrus-plugin/tree/main)

sterrors_logrus_plugin is an addition to sterrors lib to extend logrus logging formatter.

If you want to use logrus logger you can do it like that:

```go
package main

import (
	"github.com/Bibob7/sterrors-logrus-plugin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	sterrors.SetLogger(&sterrors_logrus_plugin.LogrusLogger{}) // this is not necessary, because LogrusLogger is the default logger
	err := anotherMethod()

	// second error that results from the first one
	secondErr := sterrors.E("action not possible", sterrors.SeverityError, err)

	sterrors.Log(secondErr)
}

func anotherMethod() error {
	return sterrors.E("some error message", sterrors.SeverityWarning)
}
```

Output:

```json
{
  "level": "error",
  "msg": "main.main: action not possible",
  "stack": [
    {
      "msg": "action not possible",
      "severity": "error",
      "caller": {
        "funcName": "main.main",
        "file": "/Users/root/Repositories/sterrors-logrus-formatter/examples/main.go",
        "line": 14
      }
    },
    {
      "msg": "some error message",
      "severity": "warning",
      "caller": {
        "funcName": "main.anotherMethod",
        "file": "/Users/root/Repositories/sterrors-logrus-formatter/examples/main.go",
        "line": 20
      }
    }
  ],
  "time": "2022-02-23T22:18:11+01:00"
}
```
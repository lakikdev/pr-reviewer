package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

// Logrus : implement Logger
type Logrus struct {
	*logrus.Entry

	Skip int
}

func Logger() *Logrus {
	logger := logrus.StandardLogger()

	return &Logrus{logrus.NewEntry(logger), 2}
}

func (l *Logrus) Init() *Logrus {
	l.Entry.Data["file"] = fileInfo(3)
	l.Entry.Data["func"] = fnInfo(3)
	return l
}

func (l *Logrus) WithField(key string, value interface{}) *Logrus {
	l.Entry = l.Entry.WithFields(logrus.Fields{
		key: value,
	})
	return l
}

func (l *Logrus) WithFields(fields map[string]interface{}) *Logrus {
	l.Entry = l.Entry.WithFields(fields)
	return l
}

// Info logs a message at level Info on the standard logger.
func (l *Logrus) Info(args ...interface{}) {
	if l.Entry.Data["file"] == nil {
		l.Entry.Data["file"] = fileInfo(l.Skip)
		l.Entry.Data["func"] = fnInfo((l.Skip))
	}
	l.Entry.Info(args...)
}

// Info logs a message at level Info on the standard logger.
func (l *Logrus) Error(args ...interface{}) {
	l.Entry.Error(args...)
}

// Info logs a message at level Info on the standard logger.
func (l *Logrus) ErrorWithSkip(skip int, args ...interface{}) {
	l.Entry.Data["file"] = fileInfo(l.Skip + skip)
	l.Entry.Data["func"] = fnInfo(l.Skip + skip)
	l.Entry.Info(args...)
}

func GetCallerFileName(skip int) string {
	return fileInfo(skip)
}

func GetCallerFuncName(skip int) string {
	return fnInfo(skip)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<Unknown>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/internal")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func fnInfo(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "<Unknown>"
	}

	funcName := runtime.FuncForPC(pc).Name()

	fn := funcName[strings.LastIndex(funcName, "/")+1:]

	if !strings.Contains(fn, "(") {
		return fn
	}

	fnArray := strings.Split(fn, ".")
	fn = fmt.Sprintf("%s.%s()", fnArray[0], fnArray[len(fnArray)-1])

	return fn
}

package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/errorreporting"
)

const (
	DEFAULT   = "Default"
	DEBUG     = "Debug"
	INFO      = "Info"
	NOTICE    = "Notice"
	WARNING   = "Warning"
	ERROR     = "Error"
	CRITICAL  = "Critical"
	ALERT     = "Alert"
	EMERGENCY = "Emergency"
)

type LogEntry struct {
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
	Trace    string    `json:"logging.googleapis.com/trace"`
	SpanID   string    `json:"logging.googleapis.com/spanId,omitempty"`
}

type Logger struct {
	ProjectID   string
	Trace       string
	ErrorClient *errorreporting.Client
	Request     *http.Request
}

func New(projectID string, traceID ...string) *Logger {
	trace := ""
	if len(traceID) > 0 {
		trace = fmt.Sprintf("projects/%v/traces/%v", projectID, traceID[0])
	}

	return &Logger{
		ProjectID: projectID,
		Trace:     trace,
	}
}

func (l *Logger) Log(spanID, severity, format string, a ...interface{}) {
	if err := json.NewEncoder(os.Stdout).Encode(&LogEntry{
		Time:     time.Now(),
		Trace:    l.Trace,
		SpanID:   spanID,
		Severity: severity,
		Message:  fmt.Sprintf(format, a...),
	}); err != nil {
		panic(err)
	}
}

func (l *Logger) Default(format string, a ...interface{}) {
	l.Log("", DEFAULT, format, a...)
}

func (l *Logger) DebugWith(spanID, format string, a ...interface{}) {
	l.Log(spanID, DEBUG, format, a...)
}

func (l *Logger) Debug(format string, a ...interface{}) {
	l.Log("", DEBUG, format, a...)
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.Log("", INFO, format, a...)
}

func (l *Logger) Notice(format string, a ...interface{}) {
	l.Log("", NOTICE, format, a...)
}

func (l *Logger) Warning(format string, a ...interface{}) {
	l.Log("", WARNING, format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.Log("", ERROR, format, a...)
}

func (l *Logger) Critical(format string, a ...interface{}) {
	l.Log("", CRITICAL, format, a...)
}

func (l *Logger) Alert(format string, a ...interface{}) {
	l.Log("", ALERT, format, a...)
}

func (l *Logger) Emergency(format string, a ...interface{}) {
	l.Log("", EMERGENCY, format, a...)
}

func (l *Logger) NewReport(ctx context.Context, req *http.Request) *Logger {
	c, err := errorreporting.NewClient(ctx, l.ProjectID, errorreporting.Config{})
	if err != nil {
		l.Error("new error report client: %v", err)
		return l
	}

	l.ErrorClient = c
	l.Request = req
	return l
}

func (l *Logger) ErrorAndReport(format string, a ...interface{}) {
	l.Error(format, a...)
	if l.ErrorClient == nil {
		return
	}

	for _, aa := range a {
		switch err := aa.(type) {
		case error:
			l.ErrorClient.Report(errorreporting.Entry{
				Error: err,
				Req:   l.Request,
			})
		}
	}
}

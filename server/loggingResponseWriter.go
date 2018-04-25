package server

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// LoggingResponseWriter describes methods of loggingResponseWriter helper
type LoggingResponseWriter interface {
	// WriteHeader will set and log response status code
	WriteHeader(code int)
	// WriteHeaderWithErr same as WriteHeader but also will set and log error with metadata
	WriteHeaderWithErr(code int, err error)
	// WriteJSON is a helper that will set corresponding content type, encode your data and log response
	WriteJSON(data interface{}, status ...int)
	// WriteXML is a helper that will set corresponding content type, encode your data and log response
	WriteXML(data interface{}, status ...int)
}

type loggingResponseWriter struct {
	rw         http.ResponseWriter
	remoteAddr string
	origin     string
	uri        string
	method     string
	startTime  time.Time
	statusCode int
	error      error
	logger     logrus.FieldLogger
}

func (l *loggingResponseWriter) log() {
	if l.error != nil {
		pc, fn, line, _ := runtime.Caller(3)
		l.logger.WithFields(logrus.Fields{
			"func": runtime.FuncForPC(pc).Name(),
			"file": fn,
			"line": line,
		}).Warnf("%s (%d %s) -> %s", l.uri, l.statusCode, http.StatusText(l.statusCode), l.error)
	} else {
		l.logger.Infof("%s (%d %s)", l.uri, l.statusCode, http.StatusText(l.statusCode))
	}
}

// WriteHeader will set and log response status code
func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.rw.WriteHeader(code)
	l.log()
}

// WriteHeaderWithErr same as WriteHeader but also will set and log error with metadata
func (l *loggingResponseWriter) WriteHeaderWithErr(code int, err error) {
	l.error = err
	l.statusCode = code
	l.rw.WriteHeader(code)
	l.log()
}

// WriteJSON is a helper that will set corresponding content type, encode your data and log response
func (l *loggingResponseWriter) WriteJSON(data interface{}, status ...int) {
	l.rw.Header().Set("Content-Type", "application/json")

	if len(status) != 0 {
		l.WriteHeader(status[0])
	}

	err := json.NewEncoder(l.rw).Encode(data)
	if err != nil {
		l.WriteHeaderWithErr(http.StatusInternalServerError, errors.Wrap(err, "failed to write response"))
		return
	}
	l.log()
}

// WriteXML is a helper that will set corresponding content type, encode your data and log response
func (l *loggingResponseWriter) WriteXML(data interface{}, status ...int) {
	l.rw.Header().Set("Content-Type", "application/xml")

	if len(status) != 0 {
		l.WriteHeader(status[0])
	}

	_, err := l.rw.Write([]byte(xml.Header))
	if err != nil {
		l.WriteHeaderWithErr(http.StatusInternalServerError, errors.Wrap(err, "failed to write xml header"))
		return
	}

	err = xml.NewEncoder(l.rw).Encode(data)
	if err != nil {
		l.WriteHeaderWithErr(http.StatusInternalServerError, errors.Wrap(err, "failed to write response"))
		return
	}
	l.log()
}

// ******************************************************************************************************
// We can put it in base implementation when list of errors will be stabilized and moved to separate repo
// ******************************************************************************************************
//// CheckBLError will check error from BL layer and will send corresponding http status on not nil
//// It will return false if error == nil
//func (l *loggingResponseWriter) CheckBLError(err error) bool {
//	if err == nil {
//		return false
//	}
//
//	switch err {
//	case ErrEntityNotFound:
//		l.WriteHeader(http.StatusNotFound)
//		return true
//	case ErrInvalidData:
//		l.WriteHeader(http.StatusBadRequest)
//		return true
//	case ErrDataConflict:
//		l.WriteHeader(http.StatusConflict)
//		return true
//	}
//	l.WriteHeaderWithErr(http.StatusInternalServerError, err)
//	return true
//}

func newLoggingResponseWriter(logger logrus.FieldLogger, remoteAddr, origin, uri, method string, w http.ResponseWriter) *loggingResponseWriter {
	l := &loggingResponseWriter{w, remoteAddr, origin, uri, method, time.Now(), http.StatusOK, nil, nil}
	l.logger = logger.WithFields(logrus.Fields{
		"ip":             remoteAddr,
		"origin":         origin,
		"uri":            uri,
		"method":         method,
		"executed_start": l.startTime,
		"type":           "endpoint",
	})
	return l
}

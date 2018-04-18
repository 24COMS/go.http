package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// BaseHTTPHandler stores common dependencies and methods for any http api in system.
type BaseHTTPHandler struct {
	Logger logrus.FieldLogger
}

// NewLoggingResponseWriter creates new loggingResponseWriter instance for particular request
func (h BaseHTTPHandler) NewLoggingResponseWriter(r *http.Request, w http.ResponseWriter) LoggingResponseWriter {
	return newLoggingResponseWriter(h.Logger, r.RemoteAddr, r.Referer(), r.URL.Path, r.Method, w)
}

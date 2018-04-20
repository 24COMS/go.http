package server

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// BaseHTTPHandler stores common dependencies and methods for any http api in system.
type BaseHTTPHandler struct {
	Logger logrus.FieldLogger
}

// NewLoggingResponseWriter creates new loggingResponseWriter instance for particular request
func (h BaseHTTPHandler) NewLoggingResponseWriter(w http.ResponseWriter, r *http.Request) LoggingResponseWriter {
	return newLoggingResponseWriter(h.Logger, r.RemoteAddr, r.Referer(), r.URL.Path, r.Method, w)
}

// GetLimitAndOffset will return parsed limit and offset from url.Values.
func (BaseHTTPHandler) GetLimitAndOffset(values url.Values) (offset uint64, limit uint64, err error) {
	rawOffset := values.Get("offset")
	rawLimit := values.Get("limit")

	if len(rawOffset) > 0 {
		offset, err = strconv.ParseUint(rawOffset, 10, 64)
		if err != nil {
			err = errors.Wrap(err, "failed to parse offset")
			return
		}
	}

	if len(rawLimit) > 0 {
		limit, err = strconv.ParseUint(rawLimit, 10, 64)
		if err != nil {
			err = errors.Wrap(err, "failed to parse limit")
			return
		}
	}
	return
}

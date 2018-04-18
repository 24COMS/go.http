package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/go-chi/chi/middleware"
)

const (
	swaggerBoxName  = "swagger"
	swaggerFileName = "swagger.yaml"
)

// NewBaseRouter will return new router with already registered standard endpoints
// /swagger/swagger.yaml - to serve swagger schema
// /version.json - to serve service version
// /{profilerPath}/ - to serve http profiler
func NewBaseRouter(h BaseHTTPHandler, serviceVersion, profilerPath string) (*mux.Router, error) {
	srv := mux.NewRouter()

	// SWAGGER
	box, err := rice.FindBox(swaggerBoxName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find box for "+swaggerBoxName)
	}
	file, err := box.Open(swaggerFileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open "+swaggerFileName)
	}

	APIDefinition, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read api definition")
	}

	// Register swagger handler
	srv.HandleFunc(strings.Join([]string{"", swaggerBoxName, swaggerFileName}, "/"), func(w http.ResponseWriter, r *http.Request) {
		l := h.NewLoggingResponseWriter(r, w)

		_, err := io.Copy(w, bytes.NewReader(APIDefinition))
		if err != nil {
			l.WriteHeaderWithErr(http.StatusInternalServerError, err)
		}
	})

	// VERSION
	versionResp := []byte(`{"version":"` + serviceVersion + `"}`)
	srv.HandleFunc("/version.json", func(w http.ResponseWriter, r *http.Request) {
		l := h.NewLoggingResponseWriter(r, w)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(versionResp)
		if err != nil {
			l.WriteHeaderWithErr(http.StatusInternalServerError, err)
		}
	})

	// HTTP PROFILER
	srv.Handle(profilerPath, middleware.Profiler())


	return srv, nil
}

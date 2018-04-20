package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// BaseRouterCfg stores configs for base router
type BaseRouterCfg struct {
	ServiceVersion string

	ProfilerPath string

	APIDefinition     io.Reader
	APIDefinitionPath string
}

// NewBaseRouter will return new router with already registered standard endpoints
// /version.json - to serve service version
// {cfg.ApiDefinition} - to serve api schema
// {ProfilerPath} - to serve http profiler
func NewBaseRouter(h BaseHTTPHandler, cfg BaseRouterCfg) (*mux.Router, error) {
	srv := mux.NewRouter()

	// API definition
	APIDefinition, err := ioutil.ReadAll(cfg.APIDefinition)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read api definition")
	}

	// Register swagger handler
	srv.HandleFunc(cfg.APIDefinitionPath, func(w http.ResponseWriter, r *http.Request) {
		l := h.NewLoggingResponseWriter(w, r)

		_, err := io.Copy(w, bytes.NewReader(APIDefinition))
		if err != nil {
			l.WriteHeaderWithErr(http.StatusInternalServerError, err)
		}
	})

	// VERSION
	versionResp := []byte(`{"version":"` + cfg.ServiceVersion + `"}`)
	srv.HandleFunc("/version.json", func(w http.ResponseWriter, r *http.Request) {
		l := h.NewLoggingResponseWriter(w, r)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(versionResp)
		if err != nil {
			l.WriteHeaderWithErr(http.StatusInternalServerError, err)
		}
	})

	// HTTP PROFILER
	if len(cfg.ProfilerPath) != 0 {
		srv.Handle(cfg.ProfilerPath, middleware.Profiler())
	}

	return srv, nil
}

package server

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

var cors = handlers.CORS(
	handlers.AllowCredentials(),
	handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}),
	handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"}),
)

// NewBaseServer will cinstruct http server with base configurations. Timeouts|TLS|CORS
func NewBaseServer(router http.Handler, tlsCerts []tls.Certificate, listenAddr string) *http.Server {
	return &http.Server{
		Addr:         listenAddr,
		Handler:      cors(router),
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS10,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			Certificates: tlsCerts,
		},
	}
}

package server

import (
	"fmt"
	"net/http"

	"payroll/internal/infrastructure/config"
)

func NewServer(serverConfig config.ServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:                         fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler:                      handler,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  serverConfig.ReadTimeout,
		WriteTimeout:                 serverConfig.WriteTimeout,
	}
}

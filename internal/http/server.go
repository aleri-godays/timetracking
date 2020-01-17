package http

import (
	"context"
	"fmt"
	"github.com/aleri-godays/timetracking"
	"net/http"
	"time"

	"github.com/aleri-godays/timetracking/internal/config"

	"github.com/labstack/echo/v4"
)

//Server is a http server
type Server struct {
	config     *config.Config
	httpServer *http.Server
}

//NewServer creates a new Server
func NewServer(conf *config.Config, repository timetracking.Repository) *Server {

	s := &Server{
		config: conf,
		httpServer: &http.Server{
			Addr:              fmt.Sprintf(":%d", conf.HTTPPort),
			ReadTimeout:       60 * time.Second,  // time to read request
			ReadHeaderTimeout: 10 * time.Second,  // time to read header, low value to cope with malicious behavior
			WriteTimeout:      20 * time.Second,  // time write response
			IdleTimeout:       120 * time.Second, // time between keep-alives requests before connection is closed
		},
	}
	handler := timetrackingHandler{
		repo: repository,
	}

	s.AddRoutes(handler)

	return s
}

//Start starts the echo http server
func (s *Server) Start() error {
	panic("implement me")
}

//Shutdown stops the echo http server
func (s *Server) Shutdown(ctx context.Context) error {
	panic("implement me")
}

func jsonError(c echo.Context, msg string, httpCode int) error {
	json := map[string]string{
		"request_id": c.Get("request_id").(string),
		"message":    msg,
	}
	return c.JSON(httpCode, json)
}

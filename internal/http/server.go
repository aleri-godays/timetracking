package http

import (
	"context"
	"fmt"
	"github.com/aleri-godays/timetracking"
	"net/http"
	"time"

	"github.com/aleri-godays/timetracking/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//Server is a http server
type Server struct {
	config     *config.Config
	e          *echo.Echo
	httpServer *http.Server
}

//NewServer creates a new Server
func NewServer(conf *config.Config, repository timetracking.Repository) *Server {
	e := echo.New()

	//recover from panics
	e.Use(middleware.Recover())
	//add a unique id to each request
	e.Use(middleware.RequestID())
	//add request id to the context
	e.Use(AddRequestIDToContext())
	//add a logger to the context
	e.Use(AddLoggerToContext())
	//use custom logger for all requests
	e.Use(Logger())
	//trace rest calls
	e.Use(Tracing())

	e.HideBanner = true
	e.HidePort = true

	s := &Server{
		config: conf,
		e:      e,
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
	return s.e.StartServer(s.httpServer)
}

//Shutdown stops the echo http server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

//EchoError is a wrapper for echo.NewHTTPError that contains a request id
func EchoError(c echo.Context, httpCode int, msg string) *echo.HTTPError {
	requestID := c.Get("request_id").(string)
	errMsg := fmt.Sprintf("%s. request_id: %s", msg, requestID)
	return echo.NewHTTPError(httpCode, errMsg)
}

func jsonError(c echo.Context, msg string, httpCode int) error {
	json := map[string]string{
		"request_id": c.Get("request_id").(string),
		"message":    msg,
	}
	return c.JSON(httpCode, json)
}

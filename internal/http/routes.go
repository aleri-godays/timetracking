package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//AddRoutes creates all routes
func (s *Server) AddRoutes(handler timetrackingHandler) {
	s.e.GET("/health", s.Health)

	api := s.e.Group("/api/v1")

	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(s.config.JWTSecret),
	}))

	api.POST("/timetracking", handler.Add)
	api.GET("/timetracking/:id", handler.Get)
	api.PUT("/timetracking/:id", handler.Update)
	api.DELETE("/timetracking/:id", handler.Delete)
	api.GET("/timetracking", handler.All)

	s.e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

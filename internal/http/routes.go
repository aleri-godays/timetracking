package http

//AddRoutes creates all routes
func (s *Server) AddRoutes(handler timetrackingHandler) {
	s.e.GET("/health", s.Health)

	api := s.e.Group("/api/v1")

	api.POST("/timetracking", handler.Add)
	api.GET("/timetracking/:id", handler.Get)
	api.PUT("/timetracking/:id", handler.Update)
	api.DELETE("/timetracking/:id", handler.Delete)
	api.GET("/timetracking", handler.All)
}

package http

//AddRoutes creates all routes
func (s *Server) AddRoutes(handler timetrackingHandler) {
	s.e.GET("/health", s.Health)
}

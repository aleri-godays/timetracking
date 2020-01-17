package http

import (
	"github.com/aleri-godays/timetracking"
	"github.com/labstack/echo/v4"
)

type timetrackingHandler struct {
	repo timetracking.Repository
}

func (h *timetrackingHandler) Get(c echo.Context) error {
	panic("implement me")
}

func (h *timetrackingHandler) Delete(c echo.Context) error {
	panic("implement me")
}

func (h *timetrackingHandler) Add(c echo.Context) error {
	panic("implement me")
}

func (h *timetrackingHandler) Update(c echo.Context) error {
	panic("implement me")
}

func (h *timetrackingHandler) All(c echo.Context) error {
	panic("implement me")
}

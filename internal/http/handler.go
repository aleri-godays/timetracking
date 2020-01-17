package http

import (
	"github.com/aleri-godays/timetracking"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type timetrackingHandler struct {
	repo timetracking.Repository
}

func (h *timetrackingHandler) Get(c echo.Context) error {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return jsonError(c, "id must be an integer", http.StatusBadRequest)
	}

	entry, err := h.repo.Get(c.Request().Context(), id)
	if err != nil {
		return jsonError(c, "could not get entry", http.StatusInternalServerError)
	}

	if entry == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSONPretty(http.StatusOK, entry, "  ")
}

func (h *timetrackingHandler) Delete(c echo.Context) error {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return jsonError(c, "id must be an integer", http.StatusBadRequest)
	}

	entry, err := h.repo.Get(c.Request().Context(), id)
	if err != nil {
		return jsonError(c, "could not delete entry", http.StatusInternalServerError)
	}

	if entry == nil {
		return c.NoContent(http.StatusNotFound)
	}

	err = h.repo.Delete(c.Request().Context(), id)
	if err != nil {
		return jsonError(c, "could not delete entry", http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *timetrackingHandler) Add(c echo.Context) error {
	var entry timetracking.Entry
	if err := c.Bind(&entry); err != nil {
		return jsonError(c, "invalid request body", http.StatusBadRequest)
	}

	addedEntry, err := h.repo.Add(c.Request().Context(), &entry)
	if err != nil {
		return jsonError(c, "could not save entry", http.StatusInternalServerError)
	}

	type AddResponse struct {
		ID int `json:"id"`
	}
	ar := AddResponse{ID: addedEntry.ID}

	return c.JSON(http.StatusCreated, ar)
}

func (h *timetrackingHandler) Update(c echo.Context) error {
	var entry timetracking.Entry
	if err := c.Bind(&entry); err != nil {
		return jsonError(c, "invalid request body", http.StatusBadRequest)
	}
	err := h.repo.Update(c.Request().Context(), &entry)
	if err != nil {
		return jsonError(c, "could not update entry", http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *timetrackingHandler) All(c echo.Context) error {
	entries, err := h.repo.All(c.Request().Context())
	if err != nil {
		return jsonError(c, "could not update entry", http.StatusInternalServerError)
	}

	user := c.QueryParam("user")
	if user != "" {
		filtered := make([]*timetracking.Entry, 0, 100)
		for _, entry := range entries {
			if entry.User == user {
				filtered = append(filtered, entry)
			}
		}
		entries = filtered
	}

	if len(entries) == 0 {
		return c.JSON(http.StatusOK, []string{})
	}

	return c.JSON(http.StatusOK, entries)
}

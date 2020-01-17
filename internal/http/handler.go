package http

import (
	"github.com/aleri-godays/timetracking"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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
	logger := c.Get("logger").(*log.Entry)

	var entry timetracking.Entry
	if err := c.Bind(&entry); err != nil {
		logger.WithFields(log.Fields{
			"error": err,
		}).Warn("could parse request body")
		return jsonError(c, "invalid request body", http.StatusBadRequest)
	}

	addedEntry, err := h.repo.Add(c.Request().Context(), &entry)
	if err != nil {
		logger.WithFields(log.Fields{
			"error": err,
		}).Error("could save entry")
		return jsonError(c, "could not save entry", http.StatusInternalServerError)
	}

	logger.WithFields(log.Fields{
		"entry_id": addedEntry.ID,
	}).Debug("saved entry")

	type AddResponse struct {
		ID int `json:"id"`
	}
	ar := AddResponse{ID: addedEntry.ID}

	return c.JSON(http.StatusCreated, ar)
}

func (h *timetrackingHandler) Update(c echo.Context) error {
	logger := c.Get("logger").(*log.Entry)

	var entry timetracking.Entry
	if err := c.Bind(&entry); err != nil {
		logger.WithFields(log.Fields{
			"error": err,
		}).Warn("could parse request body")
		return jsonError(c, "invalid request body", http.StatusBadRequest)
	}
	err := h.repo.Update(c.Request().Context(), &entry)
	if err != nil {
		logger.WithFields(log.Fields{
			"entry_id": entry.ID,
			"error":    err,
		}).Error("could update entry")
		return jsonError(c, "could not update entry", http.StatusInternalServerError)
	}

	logger.WithFields(log.Fields{
		"entry_id": entry.ID,
	}).Debug("updated entry")

	return c.NoContent(http.StatusNoContent)
}

func (h *timetrackingHandler) All(c echo.Context) error {
	logger := c.Get("logger").(*log.Entry)

	entries, err := h.repo.All(c.Request().Context())
	if err != nil {
		logger.WithFields(log.Fields{
			"error": err,
		}).Error("could not fetch all entries")
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

	logger.WithFields(log.Fields{
		"entry_count": len(entries),
	}).Debug("found entries")

	if len(entries) == 0 {
		return c.JSON(http.StatusOK, []string{})
	}

	return c.JSON(http.StatusOK, entries)
}

package handler

import (
	"net/http"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"
	"go-boilerplate-clean/internal/transport/apis/dto"
	logUsecase "go-boilerplate-clean/internal/usecase/notificationlogs"

	"github.com/labstack/echo/v4"
)

type NotificationLogHandler struct {
	service logUsecase.NotificationLogService
}

func NewNotificationLogHandler(service logUsecase.NotificationLogService) *NotificationLogHandler {
	return &NotificationLogHandler{service: service}
}

func (h *NotificationLogHandler) Create(c echo.Context) error {
	var req dto.CreateNotificationLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	l, err := h.service.Create(c.Request().Context(), req.ToEntity())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, l)
}

func (h *NotificationLogHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	l, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, l)
}

func (h *NotificationLogHandler) List(c echo.Context) error {
	list, err := h.service.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *NotificationLogHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateNotificationLogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	l, err := h.service.Update(c.Request().Context(), logEntity.NotificationLog{
		ID:              id,
		NotificationID:  req.NotificationID,
		UserID:          req.UserID,
		// Channel:         logEntity.Channel(req.Channel),
		SendTo:          req.SendTo,
		RenderedSubject: req.RenderedSubject,
		RenderedMessage: req.RenderedMessage,
		// Data:            req.Data,
		State:           logEntity.State(req.State),
		RetryCount:      req.RetryCount,
		ErrorMessage:   req.ErrorMessage,
		SentAt:          req.SentAt,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, l)
}

func (h *NotificationLogHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

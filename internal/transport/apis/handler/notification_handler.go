package handler

import (
	"net/http"

	"go-boilerplate-clean/internal/entity/notifications"
	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	"go-boilerplate-clean/internal/transport/apis/dto"
	notifUsecase "go-boilerplate-clean/internal/usecase/notifications"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	service notifUsecase.NotificationService
}

func NewNotificationHandler(service notifUsecase.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) Create(c echo.Context) error {
	var req dto.CreateNotificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	n, err := h.service.Create(c.Request().Context(), req.ToEntity())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, n)
}

func (h *NotificationHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	n, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, n)
}

func (h *NotificationHandler) List(c echo.Context) error {
	list, err := h.service.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *NotificationHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateNotificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	n, err := h.service.Update(c.Request().Context(), notifEntity.Notification{
		ID:                     id,
		EventKey:               req.EventKey,
		NotificationTemplateID: req.NotificationTemplateID,
		Data:                   req.Data,
		Category:               notifications.Category(req.Category),
		State:                  req.State,
		ScheduleAt:             req.ScheduleAt,
		UpdatedBy:              req.UpdatedBy,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, n)
}

func (h *NotificationHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

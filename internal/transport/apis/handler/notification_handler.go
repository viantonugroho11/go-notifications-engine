package handler

import (
	"net/http"

	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
	notifEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
	"github.com/viantonugroho11/go-notifications-engine/internal/transport/apis/dto"
	notifUsecase "github.com/viantonugroho11/go-notifications-engine/internal/usecase/notifications"

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
	param := dto.NotificationListParamFromQuery(c)
	list, err := h.service.List(c.Request().Context(), param)
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
		State:                  notifEntity.State(req.State),
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

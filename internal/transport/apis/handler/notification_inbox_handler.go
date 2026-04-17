package handler

import (
	"net/http"

	inboxEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationinbox"
	"github.com/viantonugroho11/go-notifications-engine/internal/transport/apis/dto"
	inboxUsecase "github.com/viantonugroho11/go-notifications-engine/internal/usecase/notificationinbox"

	"github.com/labstack/echo/v4"
)

type NotificationInboxHandler struct {
	service inboxUsecase.NotificationInboxService
}

func NewNotificationInboxHandler(service inboxUsecase.NotificationInboxService) *NotificationInboxHandler {
	return &NotificationInboxHandler{service: service}
}

func (h *NotificationInboxHandler) Create(c echo.Context) error {
	var req dto.CreateNotificationInboxRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	i, err := h.service.Create(c.Request().Context(), req.ToEntity())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, i)
}

func (h *NotificationInboxHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	i, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, i)
}

func (h *NotificationInboxHandler) List(c echo.Context) error {
	param := dto.NotificationInboxListParamFromQuery(c)
	list, err := h.service.List(c.Request().Context(), param)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *NotificationInboxHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateNotificationInboxRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	i, err := h.service.Update(c.Request().Context(), inboxEntity.NotificationInbox{
		ID:                id,
		UserID:            req.UserID,
		NotificationLogID: req.NotificationLogID,
		IsRead:            req.IsRead,
		ReadAt:            req.ReadAt,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, i)
}

func (h *NotificationInboxHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

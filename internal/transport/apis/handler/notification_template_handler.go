package handler

import (
	"net/http"

	tplEntity "go-boilerplate-clean/internal/entity/notificationtemplates"
	"go-boilerplate-clean/internal/transport/apis/dto"
	tplUsecase "go-boilerplate-clean/internal/usecase/notificationtemplates"

	"github.com/labstack/echo/v4"
)

type NotificationTemplateHandler struct {
	service tplUsecase.NotificationTemplateService
}

func NewNotificationTemplateHandler(service tplUsecase.NotificationTemplateService) *NotificationTemplateHandler {
	return &NotificationTemplateHandler{service: service}
}

func (h *NotificationTemplateHandler) Create(c echo.Context) error {
	var req dto.CreateNotificationTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	t, err := h.service.Create(c.Request().Context(), req.ToEntity())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, t)
}

func (h *NotificationTemplateHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	t, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (h *NotificationTemplateHandler) List(c echo.Context) error {
	list, err := h.service.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *NotificationTemplateHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateNotificationTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	t, err := h.service.Update(c.Request().Context(), tplEntity.NotificationTemplate{
		ID:            id,
		Name:          req.Name,
		Subject:       req.Subject,
		Body:          req.Body,
		PayloadSchema: req.PayloadSchema,
		Channel:       req.Channel,
		TemplateType:  req.TemplateType,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (h *NotificationTemplateHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

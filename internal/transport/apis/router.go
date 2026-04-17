package apis

import (
	"github.com/labstack/echo/v4"

	"github.com/viantonugroho11/go-notifications-engine/internal/transport/apis/handler"
	"github.com/viantonugroho11/go-notifications-engine/internal/usecase/notificationinbox"
	"github.com/viantonugroho11/go-notifications-engine/internal/usecase/notificationlogs"
	"github.com/viantonugroho11/go-notifications-engine/internal/usecase/notifications"
	"github.com/viantonugroho11/go-notifications-engine/internal/usecase/notificationtemplates"
	"github.com/viantonugroho11/go-notifications-engine/internal/usecase/users"
)

type Services struct {
	User                 users.UserService
	Notification         notifications.NotificationService
	NotificationTemplate notificationtemplates.NotificationTemplateService
	NotificationLog      notificationlogs.NotificationLogService
	NotificationInbox    notificationinbox.NotificationInboxService
}

func RegisterRoutes(e *echo.Echo, svc Services) {
	userHandler := handler.NewUserHandler(svc.User)
	notificationHandler := handler.NewNotificationHandler(svc.Notification)
	templateHandler := handler.NewNotificationTemplateHandler(svc.NotificationTemplate)
	logHandler := handler.NewNotificationLogHandler(svc.NotificationLog)
	inboxHandler := handler.NewNotificationInboxHandler(svc.NotificationInbox)

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	usersGroup := e.Group("/users")
	usersGroup.POST("", userHandler.Create)
	usersGroup.GET("", userHandler.List)
	usersGroup.GET("/:id", userHandler.GetByID)
	usersGroup.PUT("/:id", userHandler.Update)
	usersGroup.DELETE("/:id", userHandler.Delete)

	notificationsGroup := e.Group("/notifications")
	notificationsGroup.POST("", notificationHandler.Create)
	notificationsGroup.GET("", notificationHandler.List)
	notificationsGroup.GET("/:id", notificationHandler.GetByID)
	notificationsGroup.PUT("/:id", notificationHandler.Update)
	notificationsGroup.DELETE("/:id", notificationHandler.Delete)

	templatesGroup := e.Group("/notification-templates")
	templatesGroup.POST("", templateHandler.Create)
	templatesGroup.GET("", templateHandler.List)
	templatesGroup.GET("/:id", templateHandler.GetByID)
	templatesGroup.PUT("/:id", templateHandler.Update)
	templatesGroup.DELETE("/:id", templateHandler.Delete)

	logsGroup := e.Group("/notification-logs")
	logsGroup.POST("", logHandler.Create)
	logsGroup.GET("", logHandler.List)
	logsGroup.GET("/:id", logHandler.GetByID)
	logsGroup.PUT("/:id", logHandler.Update)
	logsGroup.DELETE("/:id", logHandler.Delete)

	inboxGroup := e.Group("/notification-inbox")
	inboxGroup.POST("", inboxHandler.Create)
	inboxGroup.GET("", inboxHandler.List)
	inboxGroup.GET("/:id", inboxHandler.GetByID)
	inboxGroup.PUT("/:id", inboxHandler.Update)
	inboxGroup.DELETE("/:id", inboxHandler.Delete)
}

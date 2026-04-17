package bootstrap

import (
	notifpg "github.com/viantonugroho11/go-notifications-engine/internal/repository/notification/postgres"
	inboxpg "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationinbox/postgres"
	logpg "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationlog/postgres"
	tplpg "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationtemplate/postgres"
	userpg "github.com/viantonugroho11/go-notifications-engine/internal/repository/user/postgres"
	"github.com/viantonugroho11/go-notifications-engine/internal/transport/apis"
	usecaseinbox "github.com/viantonugroho11/go-notifications-engine/internal/usecase/notificationinbox"
	usecaselog "github.com/viantonugroho11/go-notifications-engine/internal/usecase/notificationlogs"
	usecasenotif "github.com/viantonugroho11/go-notifications-engine/internal/usecase/notifications"
	usecasetpl "github.com/viantonugroho11/go-notifications-engine/internal/usecase/notificationtemplates"
	usecaseusers "github.com/viantonugroho11/go-notifications-engine/internal/usecase/users"

	"gorm.io/gorm"
)

// newServices membuat repos dan usecases untuk HTTP API.
func newServices(db *gorm.DB, publisher usecasenotif.NotificationEventPublisher) apis.Services {
	userRepo := userpg.NewUserRepository(db)
	notificationRepo := notifpg.NewNotificationRepository(db)
	templateRepo := tplpg.NewNotificationTemplateRepository(db)
	logRepo := logpg.NewNotificationLogRepository(db)
	inboxRepo := inboxpg.NewNotificationInboxRepository(db)

	return apis.Services{
		User:                 usecaseusers.NewUserService(userRepo),
		Notification:         usecasenotif.NewNotificationService(notificationRepo, templateRepo, publisher),
		NotificationTemplate: usecasetpl.NewNotificationTemplateService(templateRepo),
		NotificationLog:      usecaselog.NewNotificationLogService(logRepo),
		NotificationInbox:    usecaseinbox.NewNotificationInboxService(inboxRepo),
	}
}

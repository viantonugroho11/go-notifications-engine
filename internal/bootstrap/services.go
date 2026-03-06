package bootstrap

import (
	notifpg "go-boilerplate-clean/internal/repository/notification/postgres"
	inboxpg "go-boilerplate-clean/internal/repository/notificationinbox/postgres"
	logpg "go-boilerplate-clean/internal/repository/notificationlog/postgres"
	tplpg "go-boilerplate-clean/internal/repository/notificationtemplate/postgres"
	userpg "go-boilerplate-clean/internal/repository/user/postgres"
	"go-boilerplate-clean/internal/transport/apis"
	usecaseinbox "go-boilerplate-clean/internal/usecase/notificationinbox"
	usecaselog "go-boilerplate-clean/internal/usecase/notificationlogs"
	usecasenotif "go-boilerplate-clean/internal/usecase/notifications"
	usecasetpl "go-boilerplate-clean/internal/usecase/notificationtemplates"
	usecaseusers "go-boilerplate-clean/internal/usecase/users"

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

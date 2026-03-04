package notificationlog

import (
	"context"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"
)

type NotificationLogRepository interface {
	Create(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error)
	GetByID(ctx context.Context, id string) (logEntity.NotificationLog, error)
	List(ctx context.Context, param *logEntity.NotificationLogListParam) ([]logEntity.NotificationLog, error)
	Update(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error)
	Delete(ctx context.Context, id string) error
}

package notificationtemplate

import (
	"context"

	tplEntity "go-boilerplate-clean/internal/entity/notificationtemplates"
)

type NotificationTemplateRepository interface {
	Create(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error)
	GetByID(ctx context.Context, id string) (tplEntity.NotificationTemplate, error)
	List(ctx context.Context, param *tplEntity.NotificationTemplateListParam) ([]tplEntity.NotificationTemplate, error)
	Update(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error)
	Delete(ctx context.Context, id string) error
}

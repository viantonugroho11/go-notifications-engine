package notifications

import (
	"context"
	"strings"
	"time"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	reponotif "go-boilerplate-clean/internal/repository/notification"
	repotpl "go-boilerplate-clean/internal/repository/notificationtemplate"
	"go-boilerplate-clean/internal/shared/schema"
)

type NotificationService interface {
	Create(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	GetByID(ctx context.Context, id string) (notifEntity.Notification, error)
	List(ctx context.Context) ([]notifEntity.Notification, error)
	Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	Delete(ctx context.Context, id string) error
}

type notificationService struct {
	repo reponotif.NotificationRepository
	repoTemplate repotpl.NotificationTemplateRepository

}

func NewNotificationService(repo reponotif.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) Create(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error) {
	if err := validateNotification(n, true); err != nil {
		return notifEntity.Notification{}, err
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = time.Now()
	}
	if n.CreatedBy == "" {
		n.CreatedBy = "system"
	}
	template, err := s.repoTemplate.GetByID(ctx, n.NotificationTemplateID)
	if err != nil {
		return notifEntity.Notification{}, err
	}

	// validate payload schema
	if err := schema.ValidatePayloadSchema(template.PayloadSchema, n.Data); err != nil {
		return notifEntity.Notification{}, err
	}

	// create notification logs
	n, err = s.repo.Create(ctx, n)
	if err != nil {
		return notifEntity.Notification{}, err
	}

	return n, nil
}

func (s *notificationService) GetByID(ctx context.Context, id string) (notifEntity.Notification, error) {
	if strings.TrimSpace(id) == "" {
		return notifEntity.Notification{}, ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}

func (s *notificationService) List(ctx context.Context) ([]notifEntity.Notification, error) {
	return s.repo.List(ctx)
}

func (s *notificationService) Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error) {
	if strings.TrimSpace(n.ID) == "" {
		return notifEntity.Notification{}, ErrIDRequired
	}
	if err := validateNotification(n, false); err != nil {
		return notifEntity.Notification{}, err
	}
	now := time.Now()
	n.UpdatedAt = &now
	return s.repo.Update(ctx, n)
}

func (s *notificationService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrIDRequired
	}
	return s.repo.Delete(ctx, id)
}

func validateNotification(n notifEntity.Notification, creating bool) error {
	if strings.TrimSpace(n.EventKey) == "" {
		return ErrEventKeyRequired
	}
	if strings.TrimSpace(n.NotificationTemplateID) == "" {
		return ErrTemplateIDRequired
	}
	if strings.TrimSpace(n.Category.String()) == "" {
		return ErrCategoryRequired
	}
	if strings.TrimSpace(n.State) == "" {
		return ErrStateRequired
	}
	return nil
}



var (
	ErrIDRequired         = newValidationError("id is required")
	ErrEventKeyRequired   = newValidationError("event_key is required")
	ErrTemplateIDRequired = newValidationError("notification_template_id is required")
	ErrCategoryRequired   = newValidationError("category is required")
	ErrStateRequired      = newValidationError("state is required")
)

type validationError struct{ msg string }

func newValidationError(msg string) error { return &validationError{msg: msg} }
func (e *validationError) Error() string { return e.msg }

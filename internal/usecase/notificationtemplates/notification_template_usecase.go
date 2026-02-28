package notificationtemplates

import (
	"context"
	"strings"
	"time"

	tplEntity "go-boilerplate-clean/internal/entity/notificationtemplates"
	repotpl "go-boilerplate-clean/internal/repository/notificationtemplate"
)

type NotificationTemplateService interface {
	Create(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error)
	GetByID(ctx context.Context, id string) (tplEntity.NotificationTemplate, error)
	List(ctx context.Context) ([]tplEntity.NotificationTemplate, error)
	Update(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error)
	Delete(ctx context.Context, id string) error
}

type notificationTemplateService struct {
	repo repotpl.NotificationTemplateRepository
}

func NewNotificationTemplateService(repo repotpl.NotificationTemplateRepository) NotificationTemplateService {
	return &notificationTemplateService{repo: repo}
}

func (s *notificationTemplateService) Create(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error) {
	if err := validateTemplate(t, true); err != nil {
		return tplEntity.NotificationTemplate{}, err
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return s.repo.Create(ctx, t)
}

func (s *notificationTemplateService) GetByID(ctx context.Context, id string) (tplEntity.NotificationTemplate, error) {
	if strings.TrimSpace(id) == "" {
		return tplEntity.NotificationTemplate{}, ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}

func (s *notificationTemplateService) List(ctx context.Context) ([]tplEntity.NotificationTemplate, error) {
	return s.repo.List(ctx)
}

func (s *notificationTemplateService) Update(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error) {
	if strings.TrimSpace(t.ID) == "" {
		return tplEntity.NotificationTemplate{}, ErrIDRequired
	}
	if err := validateTemplate(t, false); err != nil {
		return tplEntity.NotificationTemplate{}, err
	}
	now := time.Now()
	t.UpdatedAt = &now
	return s.repo.Update(ctx, t)
}

func (s *notificationTemplateService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrIDRequired
	}
	return s.repo.Delete(ctx, id)
}

func validateTemplate(t tplEntity.NotificationTemplate, creating bool) error {
	if strings.TrimSpace(t.Name) == "" {
		return ErrNameRequired
	}
	if strings.TrimSpace(t.Channel) == "" {
		return ErrChannelRequired
	}
	return nil
}

var (
	ErrIDRequired      = newValidationError("id is required")
	ErrNameRequired    = newValidationError("name is required")
	ErrChannelRequired = newValidationError("channel is required")
)

type validationError struct{ msg string }

func newValidationError(msg string) error { return &validationError{msg: msg} }
func (e *validationError) Error() string { return e.msg }

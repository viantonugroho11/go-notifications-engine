package notifications

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	"go-boilerplate-clean/internal/infrastructure/broker"
	reponotif "go-boilerplate-clean/internal/repository/notification"
	repotpl "go-boilerplate-clean/internal/repository/notificationtemplate"
	"go-boilerplate-clean/internal/shared/schema"
)

type NotificationService interface {
	Create(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	GetByID(ctx context.Context, id string) (notifEntity.Notification, error)
	List(ctx context.Context, param *notifEntity.NotificationListParam) ([]notifEntity.Notification, error)
	Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	Delete(ctx context.Context, id string) error
}

type notificationService struct {
	repo         reponotif.NotificationRepository
	repoTemplate repotpl.NotificationTemplateRepository
	producer     broker.Producer
	topic        string
}

func NewNotificationService(repo reponotif.NotificationRepository, repoTemplate repotpl.NotificationTemplateRepository, producer broker.Producer, topic string) NotificationService {
	return &notificationService{
		repo:         repo,
		repoTemplate: repoTemplate,
		producer:     producer,
		topic:        topic,
	}
}

func (s *notificationService) Create(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error) {
	var (
		renderedSubject string
		renderedMessage string
	)
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
	renderedSubject = n.GenerateRenderedMessage(template.Subject)
	renderedMessage = n.GenerateRenderedMessage(template.Body)

	// create notification logs
	for i := range n.NotificationLogs {
		n.NotificationLogs[i].RenderedSubject = renderedSubject
		n.NotificationLogs[i].RenderedMessage = renderedMessage
		n.NotificationLogs[i].CreatedAt = time.Now()
	}
	n, err = s.repo.Create(ctx, n)
	if err != nil {
		return notifEntity.Notification{}, err
	}
	if err := s.publishNotificationEvent(ctx, n); err != nil {
		// log saja, jangan gagalkan create
		_ = err
	}
	return n, nil
}

func (s *notificationService) GetByID(ctx context.Context, id string) (notifEntity.Notification, error) {
	if strings.TrimSpace(id) == "" {
		return notifEntity.Notification{}, ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}

func (s *notificationService) List(ctx context.Context, param *notifEntity.NotificationListParam) ([]notifEntity.Notification, error) {
	return s.repo.List(ctx, param)
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
	updated, err := s.repo.Update(ctx, n)
	if err != nil {
		return notifEntity.Notification{}, err
	}
	if err := s.publishNotificationEvent(ctx, updated); err != nil {
		_ = err
	}
	return updated, nil
}

func (s *notificationService) publishNotificationEvent(ctx context.Context, n notifEntity.Notification) error {
	if s.producer == nil || s.topic == "" {
		return nil
	}
	msg := n.ToProducerMessage()
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, _, err = s.producer.Publish(ctx, s.topic, []byte(n.ID), payload)
	return err
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
	if strings.TrimSpace(n.State.String()) == "" {
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
func (e *validationError) Error() string  { return e.msg }

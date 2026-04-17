package notificationinbox

import (
	"context"
	"strings"
	"time"

	inboxEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationinbox"
	repoinbox "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationinbox"
)

type NotificationInboxService interface {
	Create(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error)
	GetByID(ctx context.Context, id string) (inboxEntity.NotificationInbox, error)
	List(ctx context.Context, param *inboxEntity.NotificationInboxListParam) ([]inboxEntity.NotificationInbox, error)
	Update(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error)
	Delete(ctx context.Context, id string) error
}

type notificationInboxService struct {
	repo repoinbox.NotificationInboxRepository
}

func NewNotificationInboxService(repo repoinbox.NotificationInboxRepository) NotificationInboxService {
	return &notificationInboxService{repo: repo}
}

func (s *notificationInboxService) Create(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error) {
	if err := validateInbox(i, true); err != nil {
		return inboxEntity.NotificationInbox{}, err
	}
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now()
	}
	return s.repo.Create(ctx, i)
}

func (s *notificationInboxService) GetByID(ctx context.Context, id string) (inboxEntity.NotificationInbox, error) {
	if strings.TrimSpace(id) == "" {
		return inboxEntity.NotificationInbox{}, ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}

func (s *notificationInboxService) List(ctx context.Context, param *inboxEntity.NotificationInboxListParam) ([]inboxEntity.NotificationInbox, error) {
	return s.repo.List(ctx, param)
}

func (s *notificationInboxService) Update(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error) {
	if strings.TrimSpace(i.ID) == "" {
		return inboxEntity.NotificationInbox{}, ErrIDRequired
	}
	if err := validateInbox(i, false); err != nil {
		return inboxEntity.NotificationInbox{}, err
	}
	return s.repo.Update(ctx, i)
}

func (s *notificationInboxService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrIDRequired
	}
	return s.repo.Delete(ctx, id)
}

func validateInbox(i inboxEntity.NotificationInbox, creating bool) error {
	if strings.TrimSpace(i.UserID) == "" {
		return ErrUserIDRequired
	}
	if strings.TrimSpace(i.NotificationLogID) == "" {
		return ErrNotificationLogIDRequired
	}
	return nil
}

var (
	ErrIDRequired                = newValidationError("id is required")
	ErrUserIDRequired            = newValidationError("user_id is required")
	ErrNotificationLogIDRequired = newValidationError("notification_log_id is required")
)

type validationError struct{ msg string }

func newValidationError(msg string) error { return &validationError{msg: msg} }
func (e *validationError) Error() string  { return e.msg }

package notificationlogs

import (
	"context"
	"strings"
	"time"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"
	repolog "go-boilerplate-clean/internal/repository/notificationlog"
)

type NotificationLogService interface {
	Create(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error)
	GetByID(ctx context.Context, id string) (logEntity.NotificationLog, error)
	List(ctx context.Context, param *logEntity.NotificationLogListParam) ([]logEntity.NotificationLog, error)
	Update(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error)
	Delete(ctx context.Context, id string) error
}

type notificationLogService struct {
	repo repolog.NotificationLogRepository
}

func NewNotificationLogService(repo repolog.NotificationLogRepository) NotificationLogService {
	return &notificationLogService{repo: repo}
}

func (s *notificationLogService) Create(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	if err := validateLog(l); err != nil {
		return logEntity.NotificationLog{}, err
	}
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now()
	}
	if l.State == "" {
		l.State = logEntity.StateQueued
	}
	return s.repo.Create(ctx, l)
}

func (s *notificationLogService) GetByID(ctx context.Context, id string) (logEntity.NotificationLog, error) {
	if strings.TrimSpace(id) == "" {
		return logEntity.NotificationLog{}, ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}

func (s *notificationLogService) List(ctx context.Context, param *logEntity.NotificationLogListParam) ([]logEntity.NotificationLog, error) {
	return s.repo.List(ctx, param)
}

func (s *notificationLogService) Update(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	if strings.TrimSpace(l.ID) == "" {
		return logEntity.NotificationLog{}, ErrIDRequired
	}
	if err := validateLog(l); err != nil {
		return logEntity.NotificationLog{}, err
	}
	return s.repo.Update(ctx, l)
}

func (s *notificationLogService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrIDRequired
	}
	return s.repo.Delete(ctx, id)
}

func validateLog(l logEntity.NotificationLog) error {
	if strings.TrimSpace(l.NotificationID) == "" {
		return ErrNotificationIDRequired
	}
	if strings.TrimSpace(l.UserID) == "" {
		return ErrUserIDRequired
	}
	// if strings.TrimSpace(l.Channel.String()) == "" {
	// 	return ErrChannelRequired
	// }
	if strings.TrimSpace(l.State.String()) == "" {
		return ErrStateRequired
	}
	return nil
}

var (
	ErrIDRequired           = newValidationError("id is required")
	ErrNotificationIDRequired = newValidationError("notification_id is required")
	ErrUserIDRequired       = newValidationError("user_id is required")
	ErrChannelRequired      = newValidationError("channel is required")
	ErrStateRequired        = newValidationError("state is required")
)

type validationError struct{ msg string }

func newValidationError(msg string) error { return &validationError{msg: msg} }
func (e *validationError) Error() string { return e.msg }

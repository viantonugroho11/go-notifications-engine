package notificationlogs

import (
	"context"
	"strings"
	"time"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"
	repolog "go-boilerplate-clean/internal/repository/notificationlog"
	"go-boilerplate-clean/internal/usecase/notificationlogs/states"

	"gorm.io/gorm"
)

type NotificationLogService interface {
	Create(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error)
	GetByID(ctx context.Context, id string) (logEntity.NotificationLog, error)
	List(ctx context.Context, param *logEntity.NotificationLogListParam) ([]logEntity.NotificationLog, error)
	Update(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error)
	Delete(ctx context.Context, id string) error
}

type notificationLogService struct {
	repo                repolog.NotificationLogRepository
	stateMachineFactory states.INewNotificationLogStateMachine
}

// notificationLogTransitionSaver mengimplementasikan IOnNotificationLogStateTransition untuk semua transisi (update).
type notificationLogTransitionSaver struct {
	repo repolog.NotificationLogRepository
}

func (s *notificationLogTransitionSaver) OnStateTransition(ctx context.Context, tx *gorm.DB, update logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	return s.repo.Update(ctx, update)
}

func NewNotificationLogService(repo repolog.NotificationLogRepository) NotificationLogService {
	saver := &notificationLogTransitionSaver{repo: repo}
	factory := states.NewNotificationLogStateMachineFactory(saver, saver, saver, saver, saver)
	return &notificationLogService{
		repo:                repo,
		stateMachineFactory: factory,
	}
}

func (s *notificationLogService) Create(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	if err := validateLog(l); err != nil {
		return logEntity.NotificationLog{}, err
	}
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now()
	}
	if l.State == "" {
		l.State = logEntity.StatePending
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
	current, err := s.repo.GetByID(ctx, l.ID)
	if err != nil {
		return logEntity.NotificationLog{}, err
	}
	stateMachine, err := s.stateMachineFactory.NewStateMachine(ctx, nil, &current)
	if err != nil {
		return logEntity.NotificationLog{}, err
	}
	return stateMachine.Do(ctx, nil, l)
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

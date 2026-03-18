package states

import (
	"context"
	"fmt"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	"gorm.io/gorm"
)

// INotificationState mendefinisikan perilaku per state.
type INotificationState interface {
	Do(ctx context.Context, tx *gorm.DB, update notifEntity.Notification) (notifEntity.Notification, error)
}

// INotificationStateMachine state machine untuk notification.
type INotificationStateMachine interface {
	INotificationState
	Notification() *notifEntity.Notification
}

// IOnNotificationStateTransition handler saat transisi state.
type IOnNotificationStateTransition interface {
	OnStateTransition(ctx context.Context, tx *gorm.DB, update notifEntity.Notification) (notifEntity.Notification, error)
}

// INewNotificationStateMachine factory untuk state machine.
type INewNotificationStateMachine interface {
	NewStateMachine(ctx context.Context, tx *gorm.DB, current *notifEntity.Notification) (INotificationStateMachine, error)
}

type notificationStateMachine struct {
	data    *notifEntity.Notification
	current INotificationState

	created    INotificationState
	scheduled  INotificationState
	processing INotificationState
	sent       INotificationState
	failed     INotificationState
	completed  INotificationState
}

type notificationStateMachineFactory struct {
	onToScheduled  IOnNotificationStateTransition
	onToProcessing IOnNotificationStateTransition
	onToSent       IOnNotificationStateTransition
	onToFailed     IOnNotificationStateTransition
	onToCompleted  IOnNotificationStateTransition
	onSameState    IOnNotificationStateTransition
}

// NewNotificationStateMachineFactory membuat factory state machine notification.
func NewNotificationStateMachineFactory(
	onToScheduled IOnNotificationStateTransition,
	onToProcessing IOnNotificationStateTransition,
	onToSent IOnNotificationStateTransition,
	onToFailed IOnNotificationStateTransition,
	onToCompleted IOnNotificationStateTransition,
	onSameState IOnNotificationStateTransition,
) *notificationStateMachineFactory {
	return &notificationStateMachineFactory{
		onToScheduled:  onToScheduled,
		onToProcessing: onToProcessing,
		onToSent:       onToSent,
		onToFailed:     onToFailed,
		onToCompleted:  onToCompleted,
		onSameState:    onSameState,
	}
}

// NewStateMachine membuat state machine dari notification saat ini.
func (f *notificationStateMachineFactory) NewStateMachine(ctx context.Context, tx *gorm.DB, current *notifEntity.Notification) (INotificationStateMachine, error) {
	if current == nil || current.ID == "" {
		return nil, fmt.Errorf("notification ID is required")
	}
	sm := &notificationStateMachine{data: current}

	sm.created = &createdState{sm: sm, factory: f}
	sm.scheduled = &scheduledState{sm: sm, factory: f}
	sm.processing = &processingState{sm: sm, factory: f}
	sm.sent = &sentState{sm: sm, factory: f}
	sm.failed = &failedState{sm: sm, factory: f}
	sm.completed = &completedState{sm: sm, factory: f}

	switch current.State {
	case notifEntity.StateCreated:
		sm.current = sm.created
	case notifEntity.StateScheduled:
		sm.current = sm.scheduled
	case notifEntity.StateProcessing:
		sm.current = sm.processing
	case notifEntity.StateSent:
		sm.current = sm.sent
	case notifEntity.StateFailed:
		sm.current = sm.failed
	case notifEntity.StateCompleted:
		sm.current = sm.completed
	default:
		return nil, fmt.Errorf("unknown notification state: %s", current.State)
	}
	return sm, nil
}

func (s *notificationStateMachine) Do(ctx context.Context, tx *gorm.DB, update notifEntity.Notification) (notifEntity.Notification, error) {
	return s.current.Do(ctx, tx, update)
}

func (s *notificationStateMachine) Notification() *notifEntity.Notification {
	return s.data
}

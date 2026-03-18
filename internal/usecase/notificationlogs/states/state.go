package states

import (
	"context"
	"fmt"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"

	"gorm.io/gorm"
)

// INotificationLogState mendefinisikan perilaku per state.
type INotificationLogState interface {
	Do(ctx context.Context, tx *gorm.DB, update logEntity.NotificationLog) (logEntity.NotificationLog, error)
}

// INotificationLogStateMachine state machine untuk notification log.
type INotificationLogStateMachine interface {
	INotificationLogState
	NotificationLog() *logEntity.NotificationLog
}

// IOnNotificationLogStateTransition handler saat transisi state.
type IOnNotificationLogStateTransition interface {
	OnStateTransition(ctx context.Context, tx *gorm.DB, update logEntity.NotificationLog) (logEntity.NotificationLog, error)
}

// INewNotificationLogStateMachine factory untuk state machine.
type INewNotificationLogStateMachine interface {
	NewStateMachine(ctx context.Context, tx *gorm.DB, current *logEntity.NotificationLog) (INotificationLogStateMachine, error)
}

type notificationLogStateMachine struct {
	data    *logEntity.NotificationLog
	current INotificationLogState

	pending    INotificationLogState
	processing INotificationLogState
	sent       INotificationLogState
	failed     INotificationLogState
	completed  INotificationLogState
}

type notificationLogStateMachineFactory struct {
	onToProcessing IOnNotificationLogStateTransition
	onToSent       IOnNotificationLogStateTransition
	onToFailed     IOnNotificationLogStateTransition
	onToCompleted  IOnNotificationLogStateTransition
	onSameState    IOnNotificationLogStateTransition
}

// NewNotificationLogStateMachineFactory membuat factory state machine notification log.
func NewNotificationLogStateMachineFactory(
	onToProcessing IOnNotificationLogStateTransition,
	onToSent IOnNotificationLogStateTransition,
	onToFailed IOnNotificationLogStateTransition,
	onToCompleted IOnNotificationLogStateTransition,
	onSameState IOnNotificationLogStateTransition,
) *notificationLogStateMachineFactory {
	return &notificationLogStateMachineFactory{
		onToProcessing: onToProcessing,
		onToSent:       onToSent,
		onToFailed:     onToFailed,
		onToCompleted:  onToCompleted,
		onSameState:    onSameState,
	}
}

// NewStateMachine membuat state machine dari notification log saat ini.
func (f *notificationLogStateMachineFactory) NewStateMachine(ctx context.Context, tx *gorm.DB, current *logEntity.NotificationLog) (INotificationLogStateMachine, error) {
	if current == nil || current.ID == "" {
		return nil, fmt.Errorf("notification log ID is required")
	}
	sm := &notificationLogStateMachine{data: current}

	sm.pending = &pendingState{sm: sm, factory: f}
	sm.processing = &processingState{sm: sm, factory: f}
	sm.sent = &sentState{sm: sm, factory: f}
	sm.failed = &failedState{sm: sm, factory: f}
	sm.completed = &completedState{sm: sm, factory: f}

	switch current.State {
	case logEntity.StatePending:
		sm.current = sm.pending
	case logEntity.StateProcessing:
		sm.current = sm.processing
	case logEntity.StateSent:
		sm.current = sm.sent
	case logEntity.StateFailed:
		sm.current = sm.failed
	case logEntity.StateCompleted:
		sm.current = sm.completed
	default:
		return nil, fmt.Errorf("unknown notification log state: %s", current.State)
	}
	return sm, nil
}

func (s *notificationLogStateMachine) Do(ctx context.Context, tx *gorm.DB, update logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	return s.current.Do(ctx, tx, update)
}

func (s *notificationLogStateMachine) NotificationLog() *logEntity.NotificationLog {
	return s.data
}

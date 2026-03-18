package states

import (
	"context"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"

	"gorm.io/gorm"
)

type failedState struct {
	sm      *notificationLogStateMachine
	factory *notificationLogStateMachineFactory
}

func (s *failedState) Do(ctx context.Context, tx *gorm.DB, update logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	s.sm.data = &update
	switch update.State {
	case logEntity.StateProcessing:
		return s.factory.onToProcessing.OnStateTransition(ctx, tx, update)
	default:
		return s.factory.onSameState.OnStateTransition(ctx, tx, update)
	}
}

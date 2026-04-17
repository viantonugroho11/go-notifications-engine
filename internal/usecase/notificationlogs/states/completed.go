package states

import (
	"context"

	logEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationlogs"

	"gorm.io/gorm"
)

type completedState struct {
	sm      *notificationLogStateMachine
	factory *notificationLogStateMachineFactory
}

func (s *completedState) Do(ctx context.Context, tx *gorm.DB, update logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	s.sm.data = &update
	// Completed adalah terminal state
	return s.factory.onSameState.OnStateTransition(ctx, tx, update)
}

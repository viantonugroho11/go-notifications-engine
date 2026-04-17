package states

import (
	"context"

	logEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationlogs"

	"gorm.io/gorm"
)

type sentState struct {
	sm      *notificationLogStateMachine
	factory *notificationLogStateMachineFactory
}

func (s *sentState) Do(ctx context.Context, tx *gorm.DB, update logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	s.sm.data = &update
	switch update.State {
	case logEntity.StateCompleted:
		return s.factory.onToCompleted.OnStateTransition(ctx, tx, update)
	default:
		return s.factory.onSameState.OnStateTransition(ctx, tx, update)
	}
}

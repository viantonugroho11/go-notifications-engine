package states

import (
	"context"

	logEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationlogs"

	"gorm.io/gorm"
)

type processingState struct {
	sm      *notificationLogStateMachine
	factory *notificationLogStateMachineFactory
}

func (s *processingState) Do(ctx context.Context, tx *gorm.DB, update logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	s.sm.data = &update
	switch update.State {
	case logEntity.StateSent:
		return s.factory.onToSent.OnStateTransition(ctx, tx, update)
	case logEntity.StateFailed:
		return s.factory.onToFailed.OnStateTransition(ctx, tx, update)
	default:
		return s.factory.onSameState.OnStateTransition(ctx, tx, update)
	}
}

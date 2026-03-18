package states

import (
	"context"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	"gorm.io/gorm"
)

type processingState struct {
	sm      *notificationStateMachine
	factory *notificationStateMachineFactory
}

func (s *processingState) Do(ctx context.Context, tx *gorm.DB, update notifEntity.Notification) (notifEntity.Notification, error) {
	s.sm.data = &update
	switch update.State {
	case notifEntity.StateSent:
		return s.factory.onToSent.OnStateTransition(ctx, tx, update)
	case notifEntity.StateFailed:
		return s.factory.onToFailed.OnStateTransition(ctx, tx, update)
	default:
		return s.factory.onSameState.OnStateTransition(ctx, tx, update)
	}
}

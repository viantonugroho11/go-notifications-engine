package states

import (
	"context"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	"gorm.io/gorm"
)

type createdState struct {
	sm      *notificationStateMachine
	factory *notificationStateMachineFactory
}

func (s *createdState) Do(ctx context.Context, tx *gorm.DB, update notifEntity.Notification) (notifEntity.Notification, error) {
	s.sm.data = &update
	switch update.State {
	case notifEntity.StateScheduled:
		return s.factory.onToScheduled.OnStateTransition(ctx, tx, update)
	case notifEntity.StateProcessing:
		return s.factory.onToProcessing.OnStateTransition(ctx, tx, update)
	default:
		return s.factory.onSameState.OnStateTransition(ctx, tx, update)
	}
}

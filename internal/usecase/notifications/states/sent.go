package states

import (
	"context"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	"gorm.io/gorm"
)

type sentState struct {
	sm      *notificationStateMachine
	factory *notificationStateMachineFactory
}

func (s *sentState) Do(ctx context.Context, tx *gorm.DB, update notifEntity.Notification) (notifEntity.Notification, error) {
	s.sm.data = &update
	switch update.State {
	case notifEntity.StateCompleted:
		return s.factory.onToCompleted.OnStateTransition(ctx, tx, update)
	default:
		return s.factory.onSameState.OnStateTransition(ctx, tx, update)
	}
}

package states

import (
	"context"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	"gorm.io/gorm"
)

type completedState struct {
	sm      *notificationStateMachine
	factory *notificationStateMachineFactory
}

func (s *completedState) Do(ctx context.Context, tx *gorm.DB, update notifEntity.Notification) (notifEntity.Notification, error) {
	s.sm.data = &update
	// Completed adalah terminal state; hanya allow update data tanpa ubah state
	return s.factory.onSameState.OnStateTransition(ctx, tx, update)
}

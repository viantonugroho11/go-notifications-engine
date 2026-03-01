package postgres

import (
	"context"
	"errors"

	inboxEntity "go-boilerplate-clean/internal/entity/notificationinbox"
	"go-boilerplate-clean/internal/repository/notificationinbox"
	"go-boilerplate-clean/internal/repository/notificationinbox/model"


	"gorm.io/gorm"
)

type notificationInboxRepository struct {
	db *gorm.DB
}

func NewNotificationInboxRepository(db *gorm.DB) notificationinbox.NotificationInboxRepository {
	return &notificationInboxRepository{db: db}
}

func (r *notificationInboxRepository) Create(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error) {
	m := model.ToDBNotificationInbox(i)
	err := r.db.WithContext(ctx).Create(&m).Error
	return m.ToEntity(), err
}

func (r *notificationInboxRepository) GetByID(ctx context.Context, id string) (inboxEntity.NotificationInbox, error) {
	var m model.NotificationInbox
	err := r.db.WithContext(ctx).
		Preload("NotificationLog").
		First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return inboxEntity.NotificationInbox{}, errors.New("notification inbox not found")
	}
	return m.ToEntity(), err
}

func (r *notificationInboxRepository) List(ctx context.Context) ([]inboxEntity.NotificationInbox, error) {
	var rows []model.NotificationInbox
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]inboxEntity.NotificationInbox, 0, len(rows))
	for _, m := range rows {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (r *notificationInboxRepository) Update(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error) {
	m := model.ToDBNotificationInbox(i)
	tx := r.db.WithContext(ctx).Model(&model.NotificationInbox{}).Where("id = ?", m.ID).Updates(&m)
	if tx.Error != nil {
		return inboxEntity.NotificationInbox{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return inboxEntity.NotificationInbox{}, errors.New("notification inbox not found")
	}
	return m.ToEntity(), nil
}

func (r *notificationInboxRepository) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.NotificationInbox{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("notification inbox not found")
	}
	return nil
}


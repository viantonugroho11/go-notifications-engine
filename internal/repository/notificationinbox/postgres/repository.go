package postgres

import (
	"context"
	"errors"

	inboxEntity "go-boilerplate-clean/internal/entity/notificationinbox"
	"go-boilerplate-clean/internal/repository/notificationinbox"
	"go-boilerplate-clean/internal/repository/notificationinbox/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationInboxRepository struct {
	db *gorm.DB
}

func NewNotificationInboxRepository(db *gorm.DB) notificationinbox.NotificationInboxRepository {
	return &notificationInboxRepository{db: db}
}

func (r *notificationInboxRepository) Create(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error) {
	if i.ID == "" {
		i.ID = uuid.NewString()
	}
	m := model.NotificationInbox{
		ID:                i.ID,
		UserID:            i.UserID,
		NotificationLogID: i.NotificationLogID,
		IsRead:            i.IsRead,
		ReadAt:            i.ReadAt,
		CreatedAt:         i.CreatedAt,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	return r.toEntity(m), err
}

func (r *notificationInboxRepository) GetByID(ctx context.Context, id string) (inboxEntity.NotificationInbox, error) {
	var m model.NotificationInbox
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return inboxEntity.NotificationInbox{}, errors.New("notification inbox not found")
	}
	return r.toEntity(m), err
}

func (r *notificationInboxRepository) List(ctx context.Context) ([]inboxEntity.NotificationInbox, error) {
	var rows []model.NotificationInbox
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]inboxEntity.NotificationInbox, 0, len(rows))
	for _, m := range rows {
		result = append(result, r.toEntity(m))
	}
	return result, nil
}

func (r *notificationInboxRepository) Update(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error) {
	updates := map[string]interface{}{
		"user_id":             i.UserID,
		"notification_log_id": i.NotificationLogID,
		"is_read":             i.IsRead,
		"read_at":             i.ReadAt,
	}
	tx := r.db.WithContext(ctx).Model(&model.NotificationInbox{}).Where("id = ?", i.ID).Updates(updates)
	if tx.Error != nil {
		return inboxEntity.NotificationInbox{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return inboxEntity.NotificationInbox{}, errors.New("notification inbox not found")
	}
	return i, nil
}

func (r *notificationInboxRepository) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&model.NotificationInbox{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("notification inbox not found")
	}
	return nil
}

func (r *notificationInboxRepository) toEntity(m model.NotificationInbox) inboxEntity.NotificationInbox {
	return inboxEntity.NotificationInbox{
		ID:                m.ID,
		UserID:            m.UserID,
		NotificationLogID: m.NotificationLogID,
		IsRead:            m.IsRead,
		ReadAt:            m.ReadAt,
		CreatedAt:         m.CreatedAt,
	}
}

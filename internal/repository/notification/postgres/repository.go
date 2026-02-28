package postgres

import (
	"context"
	"encoding/json"
	"errors"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	"go-boilerplate-clean/internal/repository/notification"
	"go-boilerplate-clean/internal/repository/notification/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) notification.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error) {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	dataJSON, _ := json.Marshal(n.Data)
	m := model.Notification{
		ID:                     n.ID,
		EventKey:               n.EventKey,
		NotificationTemplateID: n.NotificationTemplateID,
		Data:                   dataJSON,
		Category:               n.Category,
		State:                  n.State,
		ScheduleAt:             n.ScheduleAt,
		CreatedBy:              n.CreatedBy,
		UpdatedBy:              n.UpdatedBy,
		CreatedAt:              n.CreatedAt,
		UpdatedAt:              n.UpdatedAt,
	}
	if m.CreatedBy == "" {
		m.CreatedBy = "system"
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	return r.toEntity(m), err
}

func (r *notificationRepository) GetByID(ctx context.Context, id string) (notifEntity.Notification, error) {
	var m model.Notification
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notifEntity.Notification{}, errors.New("notification not found")
	}
	return r.toEntity(m), err
}

func (r *notificationRepository) List(ctx context.Context) ([]notifEntity.Notification, error) {
	var rows []model.Notification
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]notifEntity.Notification, 0, len(rows))
	for _, m := range rows {
		result = append(result, r.toEntity(m))
	}
	return result, nil
}

func (r *notificationRepository) Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error) {
	dataJSON, _ := json.Marshal(n.Data)
	updates := map[string]interface{}{
		"event_key":                n.EventKey,
		"notification_template_id": n.NotificationTemplateID,
		"data":                     dataJSON,
		"category":                 n.Category,
		"state":                    n.State,
		"schedule_at":              n.ScheduleAt,
		"updated_by":               n.UpdatedBy,
		"updated_at":               n.UpdatedAt,
	}
	tx := r.db.WithContext(ctx).Model(&model.Notification{}).Where("id = ?", n.ID).Updates(updates)
	if tx.Error != nil {
		return notifEntity.Notification{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return notifEntity.Notification{}, errors.New("notification not found")
	}
	return n, nil
}

func (r *notificationRepository) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&model.Notification{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("notification not found")
	}
	return nil
}

func (r *notificationRepository) toEntity(m model.Notification) notifEntity.Notification {
	var data map[string]interface{}
	if len(m.Data) > 0 {
		_ = json.Unmarshal(m.Data, &data)
	}
	return notifEntity.Notification{
		ID:                     m.ID,
		EventKey:               m.EventKey,
		NotificationTemplateID: m.NotificationTemplateID,
		Data:                   data,
		Category:               m.Category,
		State:                  m.State,
		ScheduleAt:             m.ScheduleAt,
		CreatedBy:              m.CreatedBy,
		UpdatedBy:              m.UpdatedBy,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}
}

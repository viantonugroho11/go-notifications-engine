package postgres

import (
	"context"
	"encoding/json"
	"errors"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"
	"go-boilerplate-clean/internal/repository/notificationlog"
	"go-boilerplate-clean/internal/repository/notificationlog/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationLogRepository struct {
	db *gorm.DB
}

func NewNotificationLogRepository(db *gorm.DB) notificationlog.NotificationLogRepository {
	return &notificationLogRepository{db: db}
}

func (r *notificationLogRepository) Create(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	dataJSON, _ := json.Marshal(l.Data)
	m := model.NotificationLog{
		ID:              l.ID,
		NotificationID:  l.NotificationID,
		UserID:          l.UserID,
		Channel:         l.Channel,
		SendTo:          l.SendTo,
		RenderedSubject: l.RenderedSubject,
		RenderedMessage: l.RenderedMessage,
		Data:            dataJSON,
		State:           l.State,
		RetryCount:      l.RetryCount,
		ErrorMessage:    l.ErrorMessage,
		SentAt:          l.SentAt,
		CreatedAt:       l.CreatedAt,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	return r.toEntity(m), err
}

func (r *notificationLogRepository) GetByID(ctx context.Context, id string) (logEntity.NotificationLog, error) {
	var m model.NotificationLog
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return logEntity.NotificationLog{}, errors.New("notification log not found")
	}
	return r.toEntity(m), err
}

func (r *notificationLogRepository) List(ctx context.Context) ([]logEntity.NotificationLog, error) {
	var rows []model.NotificationLog
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]logEntity.NotificationLog, 0, len(rows))
	for _, m := range rows {
		result = append(result, r.toEntity(m))
	}
	return result, nil
}

func (r *notificationLogRepository) Update(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	dataJSON, _ := json.Marshal(l.Data)
	updates := map[string]interface{}{
		"notification_id":   l.NotificationID,
		"user_id":           l.UserID,
		"channel":           l.Channel,
		"send_to":           l.SendTo,
		"rendered_subject":  l.RenderedSubject,
		"rendered_message":  l.RenderedMessage,
		"data":              dataJSON,
		"state":             l.State,
		"retry_count":       l.RetryCount,
		"error_message":     l.ErrorMessage,
		"sent_at":           l.SentAt,
	}
	tx := r.db.WithContext(ctx).Model(&model.NotificationLog{}).Where("id = ?", l.ID).Updates(updates)
	if tx.Error != nil {
		return logEntity.NotificationLog{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return logEntity.NotificationLog{}, errors.New("notification log not found")
	}
	return l, nil
}

func (r *notificationLogRepository) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&model.NotificationLog{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("notification log not found")
	}
	return nil
}

func (r *notificationLogRepository) toEntity(m model.NotificationLog) logEntity.NotificationLog {
	var data map[string]interface{}
	if len(m.Data) > 0 {
		_ = json.Unmarshal(m.Data, &data)
	}
	return logEntity.NotificationLog{
		ID:              m.ID,
		NotificationID:  m.NotificationID,
		UserID:          m.UserID,
		Channel:         m.Channel,
		SendTo:          m.SendTo,
		RenderedSubject: m.RenderedSubject,
		RenderedMessage: m.RenderedMessage,
		Data:            data,
		State:           m.State,
		RetryCount:      m.RetryCount,
		ErrorMessage:    m.ErrorMessage,
		SentAt:          m.SentAt,
		CreatedAt:       m.CreatedAt,
	}
}

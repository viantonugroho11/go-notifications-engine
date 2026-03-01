package postgres

import (
	"context"
	"errors"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"
	"go-boilerplate-clean/internal/repository/notificationlog"
	"go-boilerplate-clean/internal/repository/notificationlog/model"

	
	"gorm.io/gorm"
)

type notificationLogRepository struct {
	db *gorm.DB
}

func NewNotificationLogRepository(db *gorm.DB) notificationlog.NotificationLogRepository {
	return &notificationLogRepository{db: db}
}

func (r *notificationLogRepository) Create(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	m := model.ToDBNotificationLog(l)
	err := r.db.WithContext(ctx).Create(&m).Error
	return m.ToEntity(), err
}

func (r *notificationLogRepository) GetByID(ctx context.Context, id string) (logEntity.NotificationLog, error) {
	var m model.NotificationLog
	err := r.db.WithContext(ctx).
		First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return logEntity.NotificationLog{}, errors.New("notification log not found")
	}
	return m.ToEntity(), err
}

func (r *notificationLogRepository) List(ctx context.Context) ([]logEntity.NotificationLog, error) {
	var rows []model.NotificationLog
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]logEntity.NotificationLog, 0, len(rows))
	for _, m := range rows {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (r *notificationLogRepository) Update(ctx context.Context, l logEntity.NotificationLog) (logEntity.NotificationLog, error) {
	m := model.ToDBNotificationLog(l)
	tx := r.db.WithContext(ctx).
		Where("id = ?", m.ID).
		Updates(&m)
	if tx.Error != nil {
		return logEntity.NotificationLog{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return logEntity.NotificationLog{}, errors.New("notification log not found")
	}
	return m.ToEntity(), nil
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


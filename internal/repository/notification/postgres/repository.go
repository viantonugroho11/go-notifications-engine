package postgres

import (
	"context"
	"errors"

	notifEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
	"github.com/viantonugroho11/go-notifications-engine/internal/repository/notification"
	"github.com/viantonugroho11/go-notifications-engine/internal/repository/notification/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) notification.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error) {
	m := model.ToDBNotification(n)
	err := r.db.WithContext(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Clauses(clause.Returning{}).
		Create(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notifEntity.Notification{}, errors.New("notification not found")
	}
	if err != nil {
		return notifEntity.Notification{}, err
	}
	return m.ToEntity(), nil
}

func (r *notificationRepository) GetByID(ctx context.Context, id string) (notifEntity.Notification, error) {
	var m model.Notification
	err := r.db.WithContext(ctx).
		Preload("NotificationLogs").
		First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notifEntity.Notification{}, errors.New("notification not found")
	}
	return m.ToEntity(), err
}

func (r *notificationRepository) List(ctx context.Context, param *notifEntity.NotificationListParam) ([]notifEntity.Notification, error) {
	var rows []model.Notification
	q := model.ApplyListParam(
		r.db.WithContext(ctx).Model(&model.Notification{}).Preload("NotificationLogs"),
		param,
	)
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]notifEntity.Notification, 0, len(rows))
	for _, m := range rows {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (r *notificationRepository) Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error) {
	m := model.ToDBNotification(n)

	tx := r.db.WithContext(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Clauses(clause.Returning{}).
		Where("id = ?", m.ID).Updates(&m)
	if tx.Error != nil {
		return notifEntity.Notification{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return notifEntity.Notification{}, errors.New("notification not found")
	}
	return m.ToEntity(), nil
}

func (r *notificationRepository) Delete(ctx context.Context, id string) error {
	// Soft delete: assumes GORM SoftDelete is set up in model.Notification (e.g., with gorm.DeletedAt field)
	tx := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Notification{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("notification not found")
	}
	return nil
}

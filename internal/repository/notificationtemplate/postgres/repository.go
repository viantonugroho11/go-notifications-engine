package postgres

import (
	"context"
	"errors"

	tplEntity "go-boilerplate-clean/internal/entity/notificationtemplates"
	"go-boilerplate-clean/internal/repository/notificationtemplate"
	"go-boilerplate-clean/internal/repository/notificationtemplate/model"

	"gorm.io/gorm"
)

type notificationTemplateRepository struct {
	db *gorm.DB
}

func NewNotificationTemplateRepository(db *gorm.DB) notificationtemplate.NotificationTemplateRepository {
	return &notificationTemplateRepository{db: db}
}

func (r *notificationTemplateRepository) Create(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error) {
	m := model.ToDBNotificationTemplate(t)
	err := r.db.WithContext(ctx).Create(&m).Error
	return m.ToEntity(), err
}

func (r *notificationTemplateRepository) GetByID(ctx context.Context, id string) (tplEntity.NotificationTemplate, error) {
	var m model.NotificationTemplate
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tplEntity.NotificationTemplate{}, errors.New("notification template not found")
	}
	return m.ToEntity(), err
}

func (r *notificationTemplateRepository) List(ctx context.Context) ([]tplEntity.NotificationTemplate, error) {
	var rows []model.NotificationTemplate
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]tplEntity.NotificationTemplate, 0, len(rows))
	for _, m := range rows {
		result = append(result, m.ToEntity())
	}
	return result, nil
}

func (r *notificationTemplateRepository) Update(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error) {
	m := model.ToDBNotificationTemplate(t)
	tx := r.db.WithContext(ctx).Model(&model.NotificationTemplate{}).Where("id = ?", m.ID).Updates(&m)
	if tx.Error != nil {
		return tplEntity.NotificationTemplate{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return tplEntity.NotificationTemplate{}, errors.New("notification template not found")
	}
	return m.ToEntity(), nil
}

func (r *notificationTemplateRepository) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&model.NotificationTemplate{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("notification template not found")
	}
	return nil
}



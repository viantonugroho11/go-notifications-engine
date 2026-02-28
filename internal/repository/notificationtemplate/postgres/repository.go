package postgres

import (
	"context"
	"encoding/json"
	"errors"

	tplEntity "go-boilerplate-clean/internal/entity/notificationtemplates"
	"go-boilerplate-clean/internal/repository/notificationtemplate"
	"go-boilerplate-clean/internal/repository/notificationtemplate/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationTemplateRepository struct {
	db *gorm.DB
}

func NewNotificationTemplateRepository(db *gorm.DB) notificationtemplate.NotificationTemplateRepository {
	return &notificationTemplateRepository{db: db}
}

func (r *notificationTemplateRepository) Create(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	schemaJSON, _ := json.Marshal(t.PayloadSchema)
	m := model.NotificationTemplate{
		ID:            t.ID,
		Name:          t.Name,
		Subject:       t.Subject,
		Body:          t.Body,
		PayloadSchema: schemaJSON,
		Channel:       t.Channel,
		TemplateType:  t.TemplateType,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		DeletedAt:     t.DeletedAt,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	return r.toEntity(m), err
}

func (r *notificationTemplateRepository) GetByID(ctx context.Context, id string) (tplEntity.NotificationTemplate, error) {
	var m model.NotificationTemplate
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tplEntity.NotificationTemplate{}, errors.New("notification template not found")
	}
	return r.toEntity(m), err
}

func (r *notificationTemplateRepository) List(ctx context.Context) ([]tplEntity.NotificationTemplate, error) {
	var rows []model.NotificationTemplate
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]tplEntity.NotificationTemplate, 0, len(rows))
	for _, m := range rows {
		result = append(result, r.toEntity(m))
	}
	return result, nil
}

func (r *notificationTemplateRepository) Update(ctx context.Context, t tplEntity.NotificationTemplate) (tplEntity.NotificationTemplate, error) {
	schemaJSON, _ := json.Marshal(t.PayloadSchema)
	updates := map[string]interface{}{
		"name":            t.Name,
		"subject":         t.Subject,
		"body":            t.Body,
		"payload_schema":  schemaJSON,
		"channel":         t.Channel,
		"template_type":   t.TemplateType,
		"updated_at":      t.UpdatedAt,
	}
	tx := r.db.WithContext(ctx).Model(&model.NotificationTemplate{}).Where("id = ?", t.ID).Updates(updates)
	if tx.Error != nil {
		return tplEntity.NotificationTemplate{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return tplEntity.NotificationTemplate{}, errors.New("notification template not found")
	}
	return t, nil
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

func (r *notificationTemplateRepository) toEntity(m model.NotificationTemplate) tplEntity.NotificationTemplate {
	var schema map[string]interface{}
	if len(m.PayloadSchema) > 0 {
		_ = json.Unmarshal(m.PayloadSchema, &schema)
	}
	return tplEntity.NotificationTemplate{
		ID:            m.ID,
		Name:          m.Name,
		Subject:       m.Subject,
		Body:          m.Body,
		PayloadSchema: schema,
		Channel:       m.Channel,
		TemplateType:  m.TemplateType,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
}

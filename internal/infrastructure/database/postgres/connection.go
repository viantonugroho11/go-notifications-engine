package postgres

import (
	"context"

	notifmodel "github.com/viantonugroho11/go-notifications-engine/internal/repository/notification/model"
	inboxmodel "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationinbox/model"
	logmodel "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationlog/model"
	tplmodel "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationtemplate/model"
	usermodel "github.com/viantonugroho11/go-notifications-engine/internal/repository/user/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(ctx context.Context, dsn string) (*gorm.DB, error) {
	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}
	// Check connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&usermodel.User{},
		&notifmodel.Notification{},
		&tplmodel.NotificationTemplate{},
		&logmodel.NotificationLog{},
		&inboxmodel.NotificationInbox{},
	)
}

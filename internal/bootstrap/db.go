package bootstrap

import (
	"context"

	"github.com/viantonugroho11/go-notifications-engine/internal/config"
	pginfra "github.com/viantonugroho11/go-notifications-engine/internal/infrastructure/database/postgres"

	"gorm.io/gorm"
)

// newDB membuat koneksi DB dan menjalankan migrate.
func newDB(ctx context.Context, cfg config.Configuration) (*gorm.DB, error) {
	db, err := pginfra.Connect(ctx, cfg.PGDSN())
	if err != nil {
		return nil, err
	}
	if err := pginfra.Migrate(db); err != nil {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		return nil, err
	}
	return db, nil
}

package model

import (
	"context"

	"gorm.io/gorm"
)

func Migrate(ctx context.Context, db *gorm.DB) error {
	models := []interface{}{
		Block{},
		Transaction{},
		Log{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return err
	}
	return nil
}

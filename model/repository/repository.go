package repository

import (
	"context"
	"fmt"
	"heimdallr/conf"
	"heimdallr/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RepositoryService struct {
	DAO model.DAOAbstracter
	// TODO cache
}

func (svc *RepositoryService) Init() error {
	db, err := newDatabaseConnection()
	if err != nil {
		return err
	}
	if dao, ok := model.NewDAO(db).(*model.DAO); ok {
		svc.DAO = dao
	}

	return nil
}

func newDatabaseConnection() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.MysqlAddress), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println("Connection to mysql failed:", err)
		return nil, err
	}

	if err := model.Migrate(context.Background(), db); err != nil {
		return nil, err
	}

	return db, nil
}

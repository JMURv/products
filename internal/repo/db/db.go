package db

import (
	"fmt"
	conf "github.com/JMURv/par-pro/products/pkg/config"
	"github.com/JMURv/par-pro/products/pkg/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

func New(conf *conf.DBConfig) *Repository {
	conn, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"postgres://%s:%s@%s:%v/%s",
				conf.User,
				conf.Password,
				conf.Host,
				conf.Port,
				conf.Database,
			),
		), &gorm.Config{TranslateError: true},
	)
	if err != nil {
		zap.L().Fatal("panic occurred", zap.Any("error", err))
	}

	if err = conn.AutoMigrate(
		&model.Item{},
		&model.ItemMedia{},
		&model.ItemAttribute{},
		&model.RelatedProduct{},
		&model.Category{},
		&model.Filter{},
		&model.Promotion{},
		&model.PromotionItem{},
		&model.Favorite{},
	); err != nil {
		zap.L().Fatal("panic occurred", zap.Any("error", err))
	}
	
	return &Repository{conn: conn}
}

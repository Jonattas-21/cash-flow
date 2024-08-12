package database

import (
	"time"

	"gorm.io/gorm"
)

type IRepository interface {
	Save(iten interface{}) error
	FindByDay(date time.Time) ([]interface{}, error)
}

type Repository struct {
	Db *gorm.DB
}

func (r *Repository) Save(iten interface{}) error {
	return r.Db.Create(iten).Error
}

func (r *Repository) FindByDay(date time.Time) ([]interface{}, error) {
	var itens []interface{}
	//err := r.Db.Find(&itens).Error
	return itens, nil
}

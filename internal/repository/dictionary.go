package repository

import (
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"gorm.io/gorm"
)

type IDictionaryRepository interface {
	Create(entity *model.Dictionary) error
	CreateMany(entities []*model.Dictionary, batchSize int) error
}

type DictionaryRepository struct {
	db *gorm.DB
}

func NewDictionaryRepository(db *gorm.DB) IDictionaryRepository {
	return &DictionaryRepository{db: db}
}

func (d *DictionaryRepository) Create(entity *model.Dictionary) error {
	return d.db.Create(entity).Error
}

func (d *DictionaryRepository) CreateMany(entities []*model.Dictionary, batchSize int) error {
	return d.db.CreateInBatches(entities, batchSize).Error
}

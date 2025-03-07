package repository

import (
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"gorm.io/gorm"
)

type IDictionaryRepository interface {
	Create(entity *model.Dictionary) error
	CreateMany(entities []*model.Dictionary, batchSize int) error
	GetAll(source string, imageEmpty bool, audioEmpty bool, page int, limit int) ([]*model.Dictionary, error)
	Update(dictionary *model.Dictionary, column string, value any) error
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

func (d *DictionaryRepository) GetAll(source string, imageEmpty bool, audioEmpty bool, page int, limit int) ([]*model.Dictionary, error) {
	var dictionaries []*model.Dictionary

	// Order
	resQuery := d.db.Order("definition")

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 5
	}
	offset := (page - 1) * limit
	resQuery = resQuery.Offset(offset).Limit(limit)

	// Source
	if source != "" {
		resQuery = resQuery.Where("source = ?", source)
	}

	// Image
	if imageEmpty {
		resQuery = resQuery.Where("image LIKE ?", "https://cdn%")
	}

	// Audio
	if audioEmpty {
		resQuery = resQuery.Where("audio_gb = ? OR audio_gb IS NULL", "")
	}

	err := resQuery.Find(&dictionaries).Error

	return dictionaries, err
}

func (d *DictionaryRepository) Update(dictionary *model.Dictionary, column string, value any) error {
	return d.db.Model(&dictionary).Update(column, value).Error
}

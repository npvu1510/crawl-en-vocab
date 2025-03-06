package repository

import (
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"gorm.io/gorm"
)

type ICategoryRepository interface {
	FindAll() ([]model.Category, error)
	Create(entity *model.Category) error
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) Create(entity *model.Category) error {
	err := r.db.Create(entity).Error
	return err
}

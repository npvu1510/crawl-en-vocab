package service

import (
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"github.com/npvu1510/crawl-en-vocab/internal/repository"
)

type ICategoryService interface {
	GetAllCategories() ([]model.Category, error)
	CreateCategory(entity *model.Category) error
}

type CategoryService struct {
	repo repository.ICategoryRepository
}

func NewCategoryService(repo repository.ICategoryRepository) ICategoryService {
	return &CategoryService{repo}
}

func (c *CategoryService) GetAllCategories() ([]model.Category, error) {
	return c.repo.FindAll()
}

func (c *CategoryService) CreateCategory(entity *model.Category) error {
	return c.repo.Create(entity)
}

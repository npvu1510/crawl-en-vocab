package service

import (
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"github.com/npvu1510/crawl-en-vocab/internal/repository"
)

type IDictionaryService interface {
	CreateDictionary(entity *model.Dictionary) error
	CreateDictionaries(entities []*model.Dictionary, batchSize int) error
}

type DictionaryService struct {
	repo repository.IDictionaryRepository
}

func NewDictionaryService(repo repository.IDictionaryRepository) IDictionaryService {
	return &DictionaryService{repo}
}

func (d *DictionaryService) CreateDictionary(entity *model.Dictionary) error {
	return d.repo.Create(entity)
}
func (d *DictionaryService) CreateDictionaries(entities []*model.Dictionary, batchSize int) error {
	return d.repo.CreateMany(entities, batchSize)
}

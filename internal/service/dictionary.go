package service

import (
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"github.com/npvu1510/crawl-en-vocab/internal/repository"
)

type IDictionaryService interface {
	CreateDictionary(entity *model.Dictionary) error
	CreateDictionaries(entities []*model.Dictionary, batchSize int) error
	GetDictionaries() ([]*model.Dictionary, error)
	UpdateImage(dictionary *model.Dictionary, imgSrc string) error
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

func (d *DictionaryService) GetDictionaries() ([]*model.Dictionary, error) {
	return d.repo.GetAll()
}

func (d *DictionaryService) UpdateImage(dictionary *model.Dictionary, imgSrc string) error {
	return d.repo.Update(dictionary, "image", imgSrc)
}

package service

import (
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"github.com/npvu1510/crawl-en-vocab/internal/repository"
)

type IDictionaryService interface {
	CreateDictionary(entity *model.Dictionary) error
	CreateDictionaries(entities []*model.Dictionary, batchSize int) error
	GetDictionaries(source string, imageEmpty bool, audioEmpty bool, page int, limit int) ([]*model.Dictionary, error)
	UpdateImage(dictionary *model.Dictionary, imgSrc string) error
	UpdateAudioGb(dictionary *model.Dictionary, audioSrc string) error
	UpdateAudioUs(dictionary *model.Dictionary, audioSrc string) error
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

func (d *DictionaryService) GetDictionaries(source string, imageEmpty bool, audioEmpty bool, page int, limit int) ([]*model.Dictionary, error) {
	return d.repo.GetAll(source, imageEmpty, audioEmpty, page, limit)
}

func (d *DictionaryService) UpdateImage(dictionary *model.Dictionary, imgSrc string) error {
	return d.repo.Update(dictionary, "image", imgSrc)
}

func (d *DictionaryService) UpdateAudioGb(dictionary *model.Dictionary, audioSrc string) error {
	return d.repo.Update(dictionary, "audio_gb", audioSrc)
}

func (d *DictionaryService) UpdateAudioUs(dictionary *model.Dictionary, audioSrc string) error {
	return d.repo.Update(dictionary, "audio_us", audioSrc)
}

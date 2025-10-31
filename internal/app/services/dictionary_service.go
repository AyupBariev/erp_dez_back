package services

import (
	"erp/internal/app/models"
	"erp/internal/app/repositories"
)

type DictionaryService struct {
	dictRepo *repositories.DictionaryRepository
}

func NewDictionaryService(dictRepo *repositories.DictionaryRepository) *DictionaryService {
	return &DictionaryService{
		dictRepo: dictRepo,
	}
}

// Универсальные методы для всех типов словарей
func (s *DictionaryService) GetAll(tableName string) ([]models.BaseDictionary, error) {
	return s.dictRepo.GetAll(tableName)
}

func (s *DictionaryService) GetByID(tableName string, id int) (*models.BaseDictionary, error) {
	return s.dictRepo.GetByID(tableName, id)
}

func (s *DictionaryService) Create(tableName string, req models.CreateDictionaryRequest) (*models.BaseDictionary, error) {
	item := &models.BaseDictionary{
		Name: req.Name,
	}

	err := s.dictRepo.Create(tableName, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *DictionaryService) Update(tableName string, id int, req models.UpdateDictionaryRequest) (*models.BaseDictionary, error) {
	item, err := s.dictRepo.GetByID(tableName, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil // not found
	}

	item.Name = req.Name

	err = s.dictRepo.Update(tableName, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *DictionaryService) Delete(tableName string, id int) error {
	return s.dictRepo.Delete(tableName, id)
}

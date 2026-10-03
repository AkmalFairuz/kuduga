package service

import (
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/repository"
)

type MetaService struct {
	metaRepository *repository.MetaRepository
}

func NewMetaService(metaRepository *repository.MetaRepository) *MetaService {
	return &MetaService{metaRepository}
}

func (s *MetaService) Get(name string) (string, error) {
	return s.metaRepository.Get(name)
}

func (s *MetaService) GetAll() ([]model.Meta, error) {
	return s.metaRepository.GetAll()
}

func (s *MetaService) Set(name string, value string) error {
	return s.metaRepository.Set(name, value)
}

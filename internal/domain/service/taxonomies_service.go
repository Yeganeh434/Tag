package service

import (
	"errors"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/repository"
)

type TaxonomyService interface {
	RegisterTagRelationship(taxonomyInfo entity.Taxonomy) error
	SaveTagRelationship(ID uint64, relationshipType string) error
	GetIDByKey(key string) (uint64, error)
	GetRelatedTagsByID(ID uint64) ([]uint64, error)
	GetIDsByTitle(title string) ([]uint64, error)
}

type taxonomyService struct {
	taxonomyRepo repository.TaxonomyRepository
}

func NewTaxonomyService(taxonomyRepo repository.TaxonomyRepository) TaxonomyService {
    return &taxonomyService{taxonomyRepo: taxonomyRepo}
}

func (s *taxonomyService) RegisterTagRelationship(taxonomyInfo entity.Taxonomy) error {
	if taxonomyInfo.ID == 0 || taxonomyInfo.FromTag == 0 || taxonomyInfo.ToTag == 0 {
		return errors.New("invalid IDs")
	}
	temp:=taxonomyInfo.RelationshipType
	if temp!="" {
		if !(temp=="inclusion"||temp=="key_value"||temp=="synonym" ||temp=="antonym"){
			return errors.New("invalid relationship type")
		}
	}
	return s.taxonomyRepo.RegisterTagRelationship(taxonomyInfo)
}

func (s *taxonomyService) SaveTagRelationship(ID uint64, relationshipType string) error{
	if ID==0 {
		return errors.New("invalid ID")
	}
	if relationshipType=="" {
		return errors.New("relationship type cannot be empty")
	}
	if !(relationshipType=="inclusion"||relationshipType=="key_value"||relationshipType=="synonym" ||relationshipType=="antonym"){
		return errors.New("invalid relationship type")
	}
	
	return s.taxonomyRepo.SaveTagRelationship(ID,relationshipType)
}

func (s *taxonomyService) GetIDByKey(key string) (uint64, error){
	if key=="" {
		return 0,errors.New("key cannot be empty")
	}

	return s.taxonomyRepo.GetIDByKey(key)
}

func (s *taxonomyService) GetRelatedTagsByID(ID uint64) ([]uint64, error) {
	if ID==0 {
		return nil , errors.New("invalid ID")
	}

	return s.taxonomyRepo.GetRelatedTagsByID(ID)
}

func (s *taxonomyService) GetIDsByTitle(title string) ([]uint64, error){
	if title=="" {
		return nil,errors.New("title cannot be empty")
	}

	return s.taxonomyRepo.GetIDsByTitle(title)
}
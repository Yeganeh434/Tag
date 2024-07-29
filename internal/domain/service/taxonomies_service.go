package service

import (
	"errors"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/repository"
)

var ErrInvalidRelationshipType = errors.New("invalid relationship type")
var ErrRelationshipTypeCannotBeEmpty = errors.New("relationship type cannot be empty")
var ErrKeyCannotBeEmpty = errors.New("key cannot be empty")
var ErrNoRelationExistsWithThisID=errors.New("no relation exists with this ID")
var ErrThisRelationshipAlreadyExists=errors.New("this relationship already exists")

type TaxonomyService interface {
	RegisterTagRelationship(taxonomyInfo entity.TaxonomyEntity) error
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

func (s *taxonomyService) RegisterTagRelationship(taxonomyInfo entity.TaxonomyEntity) error {
	if taxonomyInfo.ID == 0 || taxonomyInfo.FromTag == 0 || taxonomyInfo.ToTag == 0 {
		return ErrInvalidTagID
	}
	temp := taxonomyInfo.RelationshipType
	if temp != "" {
		if !(temp == "inclusion" || temp == "key_value" || temp == "synonym" || temp == "antonym") {
			return ErrInvalidRelationshipType
		}
	}
	return s.taxonomyRepo.RegisterTagRelationship(taxonomyInfo)
}

func (s *taxonomyService) SaveTagRelationship(ID uint64, relationshipType string) error {
	if ID == 0 {
		return ErrInvalidTagID
	}
	if relationshipType == "" {
		return ErrRelationshipTypeCannotBeEmpty
	}
	if !(relationshipType == "inclusion" || relationshipType == "key_value" || relationshipType == "synonym" || relationshipType == "antonym") {
		return ErrInvalidRelationshipType
	}

	return s.taxonomyRepo.SaveTagRelationship(ID, relationshipType)
}

func (s *taxonomyService) GetIDByKey(key string) (uint64, error) {
	if key == "" {
		return 0, ErrKeyCannotBeEmpty
	}

	return s.taxonomyRepo.GetIDByKey(key)
}

func (s *taxonomyService) GetRelatedTagsByID(ID uint64) ([]uint64, error) {
	if ID == 0 {
		return nil, ErrInvalidTagID
	}

	return s.taxonomyRepo.GetRelatedTagsByID(ID)
}

func (s *taxonomyService) GetIDsByTitle(title string) ([]uint64, error) {
	if title == "" {
		return nil, ErrTitleCannotBeEmpty
	}

	return s.taxonomyRepo.GetIDsByTitle(title)
}

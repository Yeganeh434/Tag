package repository

import (
	"tag_project/internal/domain/entity"
)

type TaxonomyRepository interface {
	RegisterTagRelationship(taxonomyInfo entity.TaxonomyEntity) error
	SaveTagRelationship(ID uint64, relationshipType string) error
	GetIDByKey(key string) (uint64, error)
	GetRelatedTagsByID(ID uint64) ([]uint64, error)
	GetIDsByTitle(title string) ([]uint64, error)
}

package repository

import (
	"context"
	"service1/internal/domain/entity"
)

type TaxonomyRepository interface {
	RegisterTagRelationship(taxonomyInfo entity.TaxonomyEntity, ctx context.Context) error
	SaveTagRelationship(ID uint64, relationshipType string, ctx context.Context) error
	GetIDByKey(key string, ctx context.Context) (uint64, error)
	GetRelatedTagsByID(ID uint64, ctx context.Context) ([]uint64, error)
	GetIDsByTitle(title string, ctx context.Context) ([]uint64, error)
	GetTagsWithSameTitle(title string, ctx context.Context) ([]uint64, error)
}

package usecases

import (
	"context"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/service"
)

type TaxonomyManagementUseCase struct {
	taxonomyService service.TaxonomyService
}

func NewTaxonomyManagementUseCase(taxonomyService service.TaxonomyService) *TaxonomyManagementUseCase {
	return &TaxonomyManagementUseCase{taxonomyService: taxonomyService}
}

func (uc *TaxonomyManagementUseCase) RegisterTagRelationship(taxonomyInfo entity.TaxonomyEntity, ctx context.Context) error {
	return uc.taxonomyService.RegisterTagRelationship(taxonomyInfo,ctx)
}

func (uc *TaxonomyManagementUseCase) SaveTagRelationship(ID uint64, relationshipType string, ctx context.Context) error {
	return uc.taxonomyService.SaveTagRelationship(ID, relationshipType,ctx)
}

func (uc *TaxonomyManagementUseCase) GetIDByKey(key string, ctx context.Context) (uint64, error) {
	return uc.taxonomyService.GetIDByKey(key,ctx)
}

func (uc *TaxonomyManagementUseCase) GetRelatedTagsByID(ID uint64, ctx context.Context) ([]uint64, error) {
	return uc.taxonomyService.GetRelatedTagsByID(ID,ctx)
}

func (uc *TaxonomyManagementUseCase) GetIDsByTitle(title string, ctx context.Context) ([]uint64, error) {
	return uc.taxonomyService.GetIDsByTitle(title,ctx)
}

func (uc *TaxonomyManagementUseCase) GetTagsWithSameTitle(title string, ctx context.Context) ([]uint64, error) {
	return uc.taxonomyService.GetTagsWithSameTitle(title,ctx)
}

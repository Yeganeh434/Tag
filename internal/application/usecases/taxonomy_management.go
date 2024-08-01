package usecases

import (
    "tag_project/internal/domain/entity"
    "tag_project/internal/domain/service"
)

type TaxonomyManagementUseCase struct{
	taxonomyService service.TaxonomyService
} 

func NewTaxonomyManagementUseCase(taxonomyService service.TaxonomyService) *TaxonomyManagementUseCase {
    return &TaxonomyManagementUseCase{taxonomyService: taxonomyService}
}

func (uc *TaxonomyManagementUseCase) RegisterTagRelationship(taxonomyInfo entity.TaxonomyEntity) error{
	return uc.taxonomyService.RegisterTagRelationship(taxonomyInfo)
}

func (uc *TaxonomyManagementUseCase) SaveTagRelationship(ID uint64, relationshipType string) error{
	return uc.taxonomyService.SaveTagRelationship(ID,relationshipType)
}

func (uc *TaxonomyManagementUseCase) GetIDByKey(key string) (uint64, error){
	return uc.taxonomyService.GetIDByKey(key)
}

func (uc *TaxonomyManagementUseCase)GetRelatedTagsByID(ID uint64) ([]uint64, error){
	return uc.taxonomyService.GetRelatedTagsByID(ID)
}

func (uc *TaxonomyManagementUseCase)GetIDsByTitle(title string) ([]uint64, error){
	return uc.taxonomyService.GetIDsByTitle(title)
}

func (uc *TaxonomyManagementUseCase)GetTagsWithSameTitle(title string) ([]uint64, error){
	return uc.taxonomyService.GetTagsWithSameTitle(title)
}
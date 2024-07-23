package usecases

import (
    "tag_project/internal/domain/entity"
    "tag_project/internal/domain/service"
)

type TagManagementUseCase struct {
    tagService service.TagService
}

func NewTagManagementUseCase(tagService service.TagService) *TagManagementUseCase {
    return &TagManagementUseCase{tagService: tagService}
}

func (uc *TagManagementUseCase) RegisterTag(tag entity.Tag) error {
    return uc.tagService.RegisterTag(tag)
}

func (uc *TagManagementUseCase) UpdateTagStatus(ID uint64, isApproved string) error {
    return uc.tagService.UpdateTagStatus(ID, isApproved)
}

func (uc *TagManagementUseCase) MergeTags(originalTagID uint64, mergeTagID uint64) error {
    return uc.tagService.MergeTags(originalTagID, mergeTagID)
}

func (uc *TagManagementUseCase) IsKeyExist(key string) (bool, error) {
    return uc.tagService.IsKeyExist(key)
}

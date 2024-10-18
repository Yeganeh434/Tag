package usecases

import (
	"context"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/service"
)

type TagManagementUseCase struct {
	tagService service.TagService
}

func NewTagManagementUseCase(tagService service.TagService) *TagManagementUseCase {
	return &TagManagementUseCase{tagService: tagService}
}

func (uc *TagManagementUseCase) RegisterTag(tag entity.TagEntity, ctx context.Context) error {
	return uc.tagService.RegisterTag(tag,ctx)
}

func (uc *TagManagementUseCase) UpdateTagStatus(ID uint64, isApproved string, ctx context.Context) error {
	return uc.tagService.UpdateTagStatus(ID, isApproved,ctx)
}

func (uc *TagManagementUseCase) MergeTags(originalTagID uint64, mergeTagID uint64, ctx context.Context) error {
	return uc.tagService.MergeTags(originalTagID, mergeTagID,ctx)
}

func (uc *TagManagementUseCase) IsKeyExist(key string, ctx context.Context) (bool, error) {
	return uc.tagService.IsKeyExist(key,ctx)
}

func (uc *TagManagementUseCase) DeleteTag(ID uint64, ctx context.Context) error {
	return uc.tagService.DeleteTag(ID,ctx)
}

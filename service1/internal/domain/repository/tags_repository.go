package repository

import (
	"context"
	"service1/internal/domain/entity"
)

type TagRepository interface {
    RegisterTag(tag entity.TagEntity,ctx context.Context) error
    UpdateTagStatus(ID uint64, isApproved string,ctx context.Context) error
    MergeTags(originalTagID uint64, mergeTagID uint64,ctx context.Context) error
    IsKeyExist(key string,ctx context.Context) (bool, error)
    DeleteTag(ID uint64,ctx context.Context) (entity.TagEntity,error)
}

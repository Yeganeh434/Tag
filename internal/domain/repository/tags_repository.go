package repository

import (
    "tag_project/internal/domain/entity"
)

type TagRepository interface {
    RegisterTag(tag entity.Tag) error
    UpdateTagStatus(ID uint64, isApproved string) error
    MergeTags(originalTagID uint64, mergeTagID uint64) error
    IsKeyExist(key string) (bool, error)
}

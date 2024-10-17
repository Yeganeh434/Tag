package service

import (
	"context"
	"errors"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/repository"
)

var ErrTagKeyAlreadyExists = errors.New("tag key already exists")
var ErrTitleCannotBeEmpty = errors.New("tag title cannot be empty")
var ErrInvalidTagID = errors.New("invalid tag ID")
var ErrNoTagExistsWithThisID = errors.New("no tag exists with this ID")

type TagService interface {
	RegisterTag(tag entity.TagEntity,ctx context.Context) error
	UpdateTagStatus(ID uint64, isApproved string,ctx context.Context) error
	MergeTags(originalTagID uint64, mergeTagID uint64,ctx context.Context) error
	IsKeyExist(key string,ctx context.Context) (bool, error)
	DeleteTag(ID uint64,ctx context.Context) error
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) RegisterTag(tag entity.TagEntity,ctx context.Context) error {
	if tag.Title == "" {
		return ErrTitleCannotBeEmpty
	}
	if tag.Key != "" {
		exists, err := s.tagRepo.IsKeyExist(tag.Key,ctx)
		if err != nil {
			return err
		}
		if exists {
			return ErrTagKeyAlreadyExists
		}
	}

	return s.tagRepo.RegisterTag(tag,ctx)
}

func (s *tagService) UpdateTagStatus(ID uint64, isApproved string,ctx context.Context) error {
	if ID == 0 {
		return ErrInvalidTagID
	}

	return s.tagRepo.UpdateTagStatus(ID, isApproved,ctx)
}

func (s *tagService) MergeTags(originalTagID uint64, mergeTagID uint64,ctx context.Context) error {
	if originalTagID == 0 || mergeTagID == 0 {
		return ErrInvalidTagID
	}

	return s.tagRepo.MergeTags(originalTagID, mergeTagID,ctx)
}

func (s *tagService) IsKeyExist(key string,ctx context.Context) (bool, error) {
	if key == "" {
		return false, errors.New("key cannot be empty")
	}

	return s.tagRepo.IsKeyExist(key,ctx)
}

func (s *tagService) DeleteTag(ID uint64,ctx context.Context) error {
	if ID == 0 {
		return ErrInvalidTagID
	}

	return s.tagRepo.DeleteTag(ID,ctx)
}

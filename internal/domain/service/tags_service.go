package service

import (
	"errors"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/repository"
)

var ErrTagKeyAlreadyExists = errors.New("tag key already exists")
var ErrTitleCannotBeEmpty = errors.New("tag title cannot be empty")
var ErrInvalidTagID = errors.New("invalid tag ID")
var ErrNoTagExistsWithThisID = errors.New("no tag exists with this ID")

type TagService interface {
	RegisterTag(tag entity.TagEntity) error
	UpdateTagStatus(ID uint64, isApproved string) error
	MergeTags(originalTagID uint64, mergeTagID uint64) error
	IsKeyExist(key string) (bool, error)
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) RegisterTag(tag entity.TagEntity) error {
	if tag.Title == "" {
		return ErrTitleCannotBeEmpty
	}
	if tag.Key != "" {
		exists, err := s.tagRepo.IsKeyExist(tag.Key)
		if err != nil {
			return err
		}
		if exists {
			return ErrTagKeyAlreadyExists
		}
	}

	return s.tagRepo.RegisterTag(tag)
}

func (s *tagService) UpdateTagStatus(ID uint64, isApproved string) error {
	if ID == 0 {
		return ErrInvalidTagID
	}

	return s.tagRepo.UpdateTagStatus(ID, isApproved)
}

func (s *tagService) MergeTags(originalTagID uint64, mergeTagID uint64) error {
	if originalTagID == 0 || mergeTagID == 0 {
		return ErrInvalidTagID
	}

	return s.tagRepo.MergeTags(originalTagID, mergeTagID)
}

func (s *tagService) IsKeyExist(key string) (bool, error) {
	if key == "" {
		return false, errors.New("key cannot be empty")
	}

	return s.tagRepo.IsKeyExist(key)
}

package service

import (
    "errors"
    "tag_project/internal/domain/entity"
    "tag_project/internal/domain/repository"
)

type TagService interface {
    RegisterTag(tag entity.Tag) error
    UpdateTagStatus(ID uint64, isApproved string) error
    MergeTags(originalTagID uint64, mergeTagID uint64) error
    DoesKeyExist(key string) (bool, error)
}

type tagService struct {
    tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
    return &tagService{tagRepo: tagRepo}
}
    
func (s *tagService) RegisterTag(tag entity.Tag) error {
    if tag.Key == "" || tag.Title == "" {
        return errors.New("tag key and title cannot be empty")
    }
    exists, err := s.tagRepo.DoesKeyExist(tag.Key)
    if err != nil {
        return err
    }
    if exists {
        return errors.New("tag key already exists")
    }

    return s.tagRepo.RegisterTag(tag)
}

func (s *tagService) UpdateTagStatus(ID uint64, isApproved string) error {
    if ID == 0 {
        return errors.New("invalid tag ID")
    }
    
    return s.tagRepo.UpdateTagStatus(ID, isApproved)
}

func (s *tagService) MergeTags(originalTagID uint64, mergeTagID uint64) error {
    if originalTagID == 0 || mergeTagID == 0 {
        return errors.New("invalid tag IDs")
    }

    return s.tagRepo.MergeTags(originalTagID, mergeTagID)
}

func (s *tagService) DoesKeyExist(key string) (bool, error) {
    if key == "" {
        return false, errors.New("key cannot be empty")
    }

    return s.tagRepo.DoesKeyExist(key)
}

package mysql

import (
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/repository"
	"tag_project/internal/domain/service"

	"gorm.io/gorm"
)

type MySQLTagRepository struct {
	db *gorm.DB
}

func NewMySQLTagRepository(db *gorm.DB) repository.TagRepository {
	return &MySQLTagRepository{db: db}
}

func (r *MySQLTagRepository) RegisterTag(tagInfo entity.Tag) error {
	result := r.db.Create(&tagInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLTagRepository) UpdateTagStatus(ID uint64, isApproved string) error {
	var tag entity.Tag
	result := r.db.Model(&entity.Tag{}).Where("id=?", ID).Find(&tag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected==0 {
		return service.ErrNoTagExistsWithThisID
	}

	result = r.db.Model(&entity.Tag{}).Where("id=?", ID).Update("status", isApproved)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLTagRepository) MergeTags(originalTagID uint64, mergeTagID uint64) error {
	var tag entity.Tag
	result := r.db.Model(&entity.Tag{}).Where("id=?", originalTagID).Find(&tag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected==0 {
		return service.ErrNoTagExistsWithThisID
	}

	var firstList []entity.Taxonomy
	result = r.db.Where("from_tag=?", originalTagID).Find(&firstList)
	if result.Error != nil {
		return result.Error
	}
	for _, value := range firstList {
		value.FromTag = mergeTagID
		result = r.db.Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	var secondList []entity.Taxonomy
	result = r.db.Where("to_tag=?", originalTagID).Find(&secondList)
	if result.Error != nil {
		return result.Error
	}
	for _, value := range secondList {
		value.ToTag = mergeTagID
		result = r.db.Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *MySQLTagRepository) IsKeyExist(key string) (bool, error) {
	var tag entity.Tag
	result := r.db.Where("`key`=?", key).Find(&tag)
	if result.Error != nil {
		return true, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}


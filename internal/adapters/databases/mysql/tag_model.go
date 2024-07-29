package mysql

import "tag_project/internal/domain/entity"

type TagModel struct {
	ID          uint64 `gorm:"primary_key"`
	Title       string
	Description string
	Picture     string
	Key         string
	Status      string //approved,rejected,under_reveiw
}

func ConvertToTagModel(tag entity.Tag) TagModel {
	return TagModel{
		ID:          tag.ID,
		Title:       tag.Title,
		Description: tag.Description,
		Picture:     tag.Picture,
		Key:         tag.Key,
		Status:      tag.Status,
	}
}

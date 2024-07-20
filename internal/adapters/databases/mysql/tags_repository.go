package mysql

import "tag_project/internal/domain/entity"

func (d *Database) RegisterTag(tagInfo entity.Tag) error{
	result:=d.db.Create(&tagInfo)
	if result.Error !=nil {
		return result.Error
	}
	return nil
}

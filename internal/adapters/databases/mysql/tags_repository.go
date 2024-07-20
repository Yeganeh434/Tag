package mysql

import "tag_project/internal/domain/entity"

func (d *Database) RegisterTag(tagInfo entity.Tag) error {
	result := d.db.Create(&tagInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) UpdateTagStatus(ID uint64, isApproved string) error {
	result:=d.db.Model(&entity.Tag{}).Where("id=?",ID).Update("status",isApproved)
	if result.Error!=nil {
		return result.Error
	}
	return nil
}

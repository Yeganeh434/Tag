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
	result := d.db.Model(&entity.Tag{}).Where("id=?", ID).Update("status", isApproved)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) MergeTags(originalTagID uint64, mergeTagID uint64) error {
	var firstList []entity.Taxonomy
	result := d.db.Where("fromtag=?", originalTagID).Find(&firstList)
	if result.Error != nil {
		return result.Error
	}
	for _, value := range firstList {
		value.FromTag = mergeTagID
		result = d.db.Create(&value)
	}
	var secondList []entity.Taxonomy
	result = d.db.Where("totag=?", originalTagID).Find(&secondList)
	if result.Error != nil {
		return result.Error
	}
	for _, value := range secondList {
		value.ToTag = mergeTagID
		result = d.db.Create(&value)
	}
	return nil
}

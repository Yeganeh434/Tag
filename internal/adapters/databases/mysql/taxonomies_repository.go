package mysql

import "tag_project/internal/domain/entity"

func (d *Database) RegisterTagRelationship(taxonomyInfo entity.Taxonomy) error{
	result:=d.db.Create(&taxonomyInfo)
	if result.Error!=nil{
		return result.Error
	}
	return nil
}
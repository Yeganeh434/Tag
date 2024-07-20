package mysql

import "tag_project/internal/domain/entity"

func (d *Database) RegisterTagRelationship(taxonomyInfo entity.Taxonomy) error {
	result := d.db.Create(&taxonomyInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) SaveTagRelationship(ID uint64, relationshipType string) error {
	result:=d.db.Model(&entity.Taxonomy{}).Where("id=?",ID).Update("RelationshipType",relationshipType)
	if result.Error!=nil {
		return result.Error
	}
	return nil
}

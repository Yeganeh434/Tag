package mysql

import "tag_project/internal/domain/entity"

func (d *Database) RegisterTagRelationship(taxonomyInfo entity.Taxonomy) error {
	result := d.DB.Create(&taxonomyInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *Database) SaveTagRelationship(ID uint64, relationshipType string) error {
	result:=d.DB.Model(&entity.Taxonomy{}).Where("id=?",ID).Update("RelationshipType",relationshipType)
	if result.Error!=nil {
		return result.Error
	}
	return nil
}

func (d *Database) GetIDByKey (key string) (uint64,error) {
	var tag entity.Tag
	result:=d.DB.Where("key=?",key).First(&tag)
	if result.Error != nil {
		return 0,result.Error
	}
	if result.RowsAffected==0 {
		return 0,nil
	}
	return tag.ID,nil
}

func (d *Database) GetRelatedTagsByID (ID uint64) ([]uint64,error){
	var firstList []uint64
	result:=d.DB.Where("fromtag=?",ID).Find(&firstList)
	if result.Error != nil {
		return nil,result.Error
	}
	var secondList []uint64
	result=d.DB.Where("totag=?",ID).Find(&secondList)
	if result.Error != nil {
		return nil,result.Error
	}
	if len(firstList)==0 && len(secondList)==0 {
		return nil,nil
	}
	IDs:=append(firstList,secondList...)
	return IDs,nil
}

func (d *Database) GetIDsByTitle(title string) ([]uint64,error) {
	var IDs []uint64
	result:=d.DB.Where("title=?",title).Find(&IDs)
	if result.Error !=nil {
		return nil, result.Error
	}
	if result.RowsAffected==0 {
		return nil,nil
	}
	return IDs,nil
}

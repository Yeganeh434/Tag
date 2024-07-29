package mysql

import (
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/repository"
	"tag_project/internal/domain/service"

	"gorm.io/gorm"
)

type MySQLTaxonomyRepository struct {
	db *gorm.DB
}

func NewMySQLTaxonomyRepository(db *gorm.DB) repository.TaxonomyRepository {
	return &MySQLTaxonomyRepository{db: db}
}

func (r *MySQLTaxonomyRepository) RegisterTagRelationship(taxonomyInfo entity.TaxonomyEntity) error {
	taxonomyModel:=ConvertToTaxonomyModel(taxonomyInfo)
	var tag1 Tag
	result := r.db.Model(&Tag{}).Where("ID=?", taxonomyModel.FromTag).Find(&tag1)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNoTagExistsWithThisID
	}

	var tag2 Tag
	result = r.db.Model(&Tag{}).Where("ID=?", taxonomyModel.ToTag).Find(&tag2)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNoTagExistsWithThisID
	}

	var taxonomy Taxonomy
	result=r.db.Model(&Taxonomy{}).Where("from_tag=? AND to_tag=?",taxonomyModel.FromTag,taxonomyModel.ToTag).Find(&taxonomy)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return service.ErrThisRelationshipAlreadyExists
	}

	result = r.db.Create(&taxonomyModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLTaxonomyRepository) SaveTagRelationship(ID uint64, relationshipType string) error {
	var taxonomy Taxonomy
	result := r.db.Model(&Taxonomy{}).Where("ID=?", ID).Find(&taxonomy)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNoRelationExistsWithThisID
	}

	result = r.db.Model(&Taxonomy{}).Where("id=?", ID).Update("RelationshipType", relationshipType)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLTaxonomyRepository) GetIDByKey(key string) (uint64, error) {
	var tag Tag
	result := r.db.Where("`key`=?", key).Find(&tag)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, nil
	}
	return tag.ID, nil
}

func (r *MySQLTaxonomyRepository) GetRelatedTagsByID(ID uint64) ([]uint64, error) {
	var tag Tag
	result := r.db.Model(&Tag{}).Where("id=?", ID).Find(&tag)
	if result.Error != nil {
		return nil,result.Error
	}
	if result.RowsAffected==0 {
		return nil,service.ErrNoTagExistsWithThisID
	}

	var firstTaxonomy []Taxonomy
	result = r.db.Where("from_tag=?", ID).Find(&firstTaxonomy)
	if result.Error != nil {
		return nil, result.Error
	}

	var secondTaxonomy []Taxonomy
	result = r.db.Where("to_tag=?", ID).Find(&secondTaxonomy)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(firstTaxonomy) == 0 && len(secondTaxonomy) == 0 {
		return nil, nil
	}
	var IDs []uint64
	for _, value := range firstTaxonomy {
		IDs = append(IDs, value.ToTag)
	}
	for _, value := range secondTaxonomy {
		IDs = append(IDs, value.FromTag)
	}
	return IDs, nil
}

func (r *MySQLTaxonomyRepository) GetIDsByTitle(title string) ([]uint64, error) {
	var tags []Tag
	result := r.db.Where("title=?", title).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	var IDs []uint64
	for _,value :=range tags {
		IDs=append(IDs, value.ID)
	}
	return IDs, nil
}

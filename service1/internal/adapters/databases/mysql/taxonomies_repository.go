package mysql

import (
	"context"
	"service1/internal/config"
	"service1/internal/domain/entity"
	"service1/internal/domain/repository"
	"service1/internal/domain/service"

	"gorm.io/gorm"
)

type MySQLTaxonomyRepository struct {
	db *gorm.DB
}

func NewMySQLTaxonomyRepository(db *gorm.DB) repository.TaxonomyRepository {
	return &MySQLTaxonomyRepository{db: db}
}

func (r *MySQLTaxonomyRepository) RegisterTagRelationship(taxonomyInfo entity.TaxonomyEntity, ctx context.Context) error {
	ctx, span := config.Tracer.Start(ctx, "RegisterTagRelationship_database")
	defer span.End()

	taxonomyModel := ConvertToTaxonomyModel(taxonomyInfo)
	var tag1 Tag
	result := r.db.WithContext(ctx).Model(&Tag{}).Where("ID=?", taxonomyModel.FromTag).Find(&tag1)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNoTagExistsWithThisID
	}

	var tag2 Tag
	result = r.db.WithContext(ctx).Model(&Tag{}).Where("ID=?", taxonomyModel.ToTag).Find(&tag2)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNoTagExistsWithThisID
	}

	var taxonomy Taxonomy
	result = r.db.WithContext(ctx).Model(&Taxonomy{}).Where("from_tag=? AND to_tag=?", taxonomyModel.FromTag, taxonomyModel.ToTag).Find(&taxonomy)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return service.ErrThisRelationshipAlreadyExists
	}

	result = r.db.WithContext(ctx).Create(&taxonomyModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLTaxonomyRepository) SaveTagRelationship(ID uint64, relationshipType string,ctx context.Context) error {
	ctx, span := config.Tracer.Start(ctx, "SaveTagRelationship_database")
	defer span.End()

	var taxonomy Taxonomy
	result := r.db.WithContext(ctx).Model(&Taxonomy{}).Where("ID=?", ID).Find(&taxonomy)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNoRelationExistsWithThisID
	}

	result = r.db.WithContext(ctx).Model(&Taxonomy{}).Where("id=?", ID).Update("RelationshipType", relationshipType)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLTaxonomyRepository) GetIDByKey(key string,ctx context.Context) (uint64, error) {
	ctx, span := config.Tracer.Start(ctx, "GetIDByKey_database")
	defer span.End()

	var tag Tag
	result := r.db.WithContext(ctx).Where("`key`=?", key).Find(&tag)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, nil
	}
	return tag.ID, nil
}

func (r *MySQLTaxonomyRepository) GetRelatedTagsByID(ID uint64,ctx context.Context) ([]uint64, error) {
	ctx, span := config.Tracer.Start(ctx, "GetRelatedTagsByID_database")
	defer span.End()

	var tag Tag
	result := r.db.WithContext(ctx).Model(&Tag{}).Where("id=?", ID).Find(&tag)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, service.ErrNoTagExistsWithThisID
	}

	var firstTaxonomy []Taxonomy
	result = r.db.WithContext(ctx).Where("from_tag=?", ID).Find(&firstTaxonomy)
	if result.Error != nil {
		return nil, result.Error
	}

	var secondTaxonomy []Taxonomy
	result = r.db.WithContext(ctx).Where("to_tag=?", ID).Find(&secondTaxonomy)
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

func (r *MySQLTaxonomyRepository) GetIDsByTitle(title string,ctx context.Context) ([]uint64, error) {
	ctx, span := config.Tracer.Start(ctx, "GetIDsByTitle_database")
	defer span.End()

	var tags []Tag
	result := r.db.WithContext(ctx).Where("title=?", title).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	var IDs []uint64
	for _, value := range tags {
		IDs = append(IDs, value.ID)
	}
	return IDs, nil
}

func (r *MySQLTaxonomyRepository) GetTagsWithSameTitle(title string,ctx context.Context) ([]uint64, error) {
	ctx, span := config.Tracer.Start(ctx, "GetTagsWithSameTitle_database")
	defer span.End()

	var categoryTag Tag
	result := r.db.WithContext(ctx).Where("title=?", "categories").Find(&categoryTag)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	var parentTag Tag
	parentTitle := "MainNode_" + title
	result = r.db.WithContext(ctx).Where("title=?", parentTitle).Find(&parentTag)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	var childsTag []Taxonomy
	result = r.db.WithContext(ctx).Where("from_tag=?", parentTag.ID).Find(&childsTag)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	var IDs []uint64
	for _, value := range childsTag {
		IDs = append(IDs, value.ID)
	}
	return IDs, nil
}

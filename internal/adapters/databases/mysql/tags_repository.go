package mysql

import (
	"context"
	"tag_project/internal/application/usecases"
	"tag_project/internal/config"
	"tag_project/internal/domain/entity"
	"tag_project/internal/domain/repository"
	"tag_project/internal/domain/service"

	// "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type MySQLTagRepository struct {
	db *gorm.DB
}

func NewMySQLTagRepository(db *gorm.DB) repository.TagRepository {
	return &MySQLTagRepository{db: db}
}

func (r *MySQLTagRepository) RegisterTag(tagInfo entity.TagEntity, ctx context.Context) error {
	ctx, span := config.Tracer.Start(ctx, "RegisterTag_database")
	defer span.End()

	tagModel := ConvertToTagModel(tagInfo)

	var categoryTag Tag
	result := r.db.WithContext(ctx).Where("title=?", "categories").Find(&categoryTag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		tagID, err := usecases.GenerateID()
		if err != nil {
			return err
		}
		categoryTag = Tag{
			ID:     tagID,
			Title:  "categories",
			Key:    "categories",
			Status: "approved",
		}
		result = r.db.Create(&categoryTag)
		if result.Error != nil {
			return result.Error
		}
	}

	var parentTag Tag
	parentTitle := "MainNode_" + tagModel.Title
	result = r.db.WithContext(ctx).Where("title=?", parentTitle).Find(&parentTag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		tagID, err := usecases.GenerateID()
		if err != nil {
			return err
		}
		parentTag = Tag{
			ID:     tagID,
			Title:  parentTitle,
			Status: "approved",
		}
		result = r.db.WithContext(ctx).Create(&parentTag)
		if result.Error != nil {
			return result.Error
		}

		taxonomyID, err := usecases.GenerateID()
		if err != nil {
			return err
		}
		taxonomy := Taxonomy{
			ID:               taxonomyID,
			FromTag:          categoryTag.ID,
			ToTag:            parentTag.ID,
			RelationshipType: "inclusion",
			Status:           "active",
		}
		result = r.db.WithContext(ctx).Create(&taxonomy)
		if result.Error != nil {
			return result.Error
		}
	}

	//register new tag
	result = r.db.WithContext(ctx).Create(&tagModel)
	if result.Error != nil {
		return result.Error
	}

	taxonomyID, err := usecases.GenerateID()
	if err != nil {
		return err
	}
	taxonomy := Taxonomy{
		ID:               taxonomyID,
		FromTag:          parentTag.ID,
		ToTag:            tagModel.ID,
		RelationshipType: "inclusion",
		Status:           "active",
	}
	result = r.db.WithContext(ctx).Create(&taxonomy)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MySQLTagRepository) UpdateTagStatus(ID uint64, isApproved string,ctx context.Context) error {
	ctx,span:=config.Tracer.Start(ctx,"UpdateTagStatus_database")
	defer span.End()

	var tag Tag
	result := r.db.WithContext(ctx).Model(&Tag{}).Where("id=?", ID).Find(&tag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNoTagExistsWithThisID
	}

	result = r.db.WithContext(ctx).Model(&Tag{}).Where("id=?", ID).Update("status", isApproved)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *MySQLTagRepository) MergeTags(originalTagID uint64, mergeTagID uint64,ctx context.Context) error {
	ctx,span:=config.Tracer.Start(ctx,"MergeTags_database")
	span.End()

	var tag Tag
	result := r.db.WithContext(ctx).Model(&Tag{}).Where("id=?", originalTagID).Find(&tag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return service.ErrNoTagExistsWithThisID
	}

	var firstList []Taxonomy
	result = r.db.WithContext(ctx).Where("from_tag=?", originalTagID).Find(&firstList)
	if result.Error != nil {
		return result.Error
	}
	for _, value := range firstList {
		value.FromTag = mergeTagID
		result = r.db.WithContext(ctx).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	var secondList []Taxonomy
	result = r.db.WithContext(ctx).Where("to_tag=?", originalTagID).Find(&secondList)
	if result.Error != nil {
		return result.Error
	}
	for _, value := range secondList {
		value.ToTag = mergeTagID
		result = r.db.WithContext(ctx).Create(&value)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r *MySQLTagRepository) IsKeyExist(key string,ctx context.Context) (bool, error) {
	ctx,span:=config.Tracer.Start(ctx,"IsKeyExist_database")
	span.End()

	var tag Tag
	result := r.db.WithContext(ctx).Where("`key`=?", key).Find(&tag)
	if result.Error != nil {
		return true, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

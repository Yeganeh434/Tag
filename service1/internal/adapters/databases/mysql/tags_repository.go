package mysql

import (
	"context"
	"service1/internal/application/usecases"
	"service1/internal/config"
	"service1/internal/domain/entity"
	"service1/internal/domain/repository"
	"service1/internal/domain/service"

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

func (r *MySQLTagRepository) UpdateTagStatus(ID uint64, isApproved string, ctx context.Context) error {
	ctx, span := config.Tracer.Start(ctx, "UpdateTagStatus_database")
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

func (r *MySQLTagRepository) MergeTags(originalTagID uint64, mergeTagID uint64, ctx context.Context) error {
	ctx, span := config.Tracer.Start(ctx, "MergeTags_database")
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

func (r *MySQLTagRepository) IsKeyExist(key string, ctx context.Context) (bool, error) {
	ctx, span := config.Tracer.Start(ctx, "IsKeyExist_database")
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

func (r *MySQLTagRepository) DeleteTag(ID uint64, ctx context.Context) (entity.TagEntity, error) {
	ctx, span := config.Tracer.Start(ctx, "DeleteTag_database")
	span.End()

	var tag Tag
	result := r.db.WithContext(ctx).Where("id=?", ID).Find(&tag)
	tagEntity := entity.TagEntity{
		ID:          tag.ID,
		Title:       tag.Title,
		Description: tag.Description,
		Picture:     tag.Picture,
		Key:         tag.Key,
		Status:      tag.Status,
	}
	if result.Error != nil {
		return tagEntity, result.Error
	}
	if result.RowsAffected == 0 {
		return tagEntity, service.ErrNoTagExistsWithThisID
	}

	var parentTag Tag
	title := "MainNode_" + tag.Title
	result = r.db.WithContext(ctx).Where("title=?", title).Find(&parentTag)
	if result.Error != nil {
		return tagEntity, result.Error
	}

	var count int64
	result = r.db.WithContext(ctx).Model(&Taxonomy{}).Where("from_tag=?", parentTag.ID).Count(&count)
	if result.Error != nil {
		return tagEntity, result.Error
	}

	if count == 1 {
		result = r.db.WithContext(ctx).Delete(&Tag{}, parentTag.ID)
		if result.Error != nil {
			return tagEntity, result.Error
		}
		result = r.db.WithContext(ctx).Where("to_tag=?", parentTag.ID).Delete(&Taxonomy{})
		if result.Error != nil {
			return tagEntity, result.Error
		}
	}

	result = r.db.WithContext(ctx).Delete(&Tag{}, ID)
	if result.Error != nil {
		return tagEntity, result.Error
	}

	result = r.db.WithContext(ctx).Where("from_tag=?", ID).Delete(&Taxonomy{})
	if result.Error != nil {
		return tagEntity, result.Error
	}

	result = r.db.WithContext(ctx).Where("to_tag=?", ID).Delete(&Taxonomy{})
	if result.Error != nil {
		return tagEntity, result.Error
	}

	return tagEntity, nil
}

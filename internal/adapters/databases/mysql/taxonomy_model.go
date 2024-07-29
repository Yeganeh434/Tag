package mysql

import "tag_project/internal/domain/entity"

type TaxonomyModel struct {
	ID               uint64 `gorm:"primary_key"`
	FromTag          uint64
	ToTag            uint64
	RelationshipType string //inclusion,key_value,synonym,antonym
	Status           string
}

func ConvertToTaxonomyModel(taxonomy entity.Taxonomy) TaxonomyModel {
	return TaxonomyModel{
		ID:               taxonomy.ID,
		FromTag:          taxonomy.FromTag,
		ToTag:            taxonomy.ToTag,
		RelationshipType: taxonomy.RelationshipType,
		Status:           taxonomy.Status,
	}
}

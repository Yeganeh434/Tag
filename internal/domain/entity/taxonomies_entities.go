package entity

type TaxonomyEntity struct {
	ID               uint64
	FromTag          uint64
	ToTag            uint64
	RelationshipType string   //inclusion,key_value,synonym,antonym
	Status           string
}

package entity

type Taxonomy struct {
	ID               uint64 `gorm:"primary_key"`
	FromTag          uint64
	ToTag            uint64
	RelationshipType string   //inclusion,key_value,synonym,antonym
	Status           string
}

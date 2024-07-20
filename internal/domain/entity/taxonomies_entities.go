package entity

type Taxonomy struct {
	ID               uint64 `gorm:"primary_key"`
	FromTagID        uint64
	ToTagID          uint64
	RelationshipType string
	Status           string
}

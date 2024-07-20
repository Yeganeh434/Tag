package entity

type Tag struct {
	ID          uint64 `gorm:"primary_key"`
	Title       string
	Description string
	Picture     string
	Key         string
	Status      string
}

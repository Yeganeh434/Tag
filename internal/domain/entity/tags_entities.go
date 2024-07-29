package entity

type Tag struct {
	ID          uint64
	Title       string
	Description string
	Picture     string
	Key         string
	Status      string //approved,rejected,under_reveiw
}



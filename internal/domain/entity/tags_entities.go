package entity

type TagEntity struct {
	ID          uint64
	Title       string
	Description string
	Picture     string
	Key         string
	Status      string //approved,rejected,under_reveiw
}



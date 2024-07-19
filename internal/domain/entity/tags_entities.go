package entity

type Tag struct {
	Title       string
	Description string
	Picture     string 
	Key         string
	Status      string
}

type taxonomy struct {
	FromTag          string
	ToTag            string
	RelationshipType string
	status           string
}

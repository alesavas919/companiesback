package models

type User struct {
	ID       int8   `json:"id"` //JSON ASSIGN AT THE DB STRUCTURE uuid.UUID user.ID = uuid.New()
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // "-" SKIP THE ATTRIBUTE RESPONSE
}

// ASSIGN THE TABLE NAME INTO DATABASE
func (User) TableName() string {
	return "users"
}

package models

type User struct {
	ID       uint   `gorm:"autoIncrement"`
	Name     string `gorm:"not null" json:"name"`
	Token    string
	Keystore string
}

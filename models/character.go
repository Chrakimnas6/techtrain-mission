package models

type Character struct {
	ID   uint   `gorm:"autoIncrement"`
	Name string `gorm:"not null" json:"name"`
}

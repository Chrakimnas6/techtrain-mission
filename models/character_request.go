package models

type CharacterRequest struct {
	ID   uint    `gorm:"autoIncrement"`
	Name string  `gorm:"not null" json:"name"`
	Rank string  `gorm:"not null" json:"rank"`
	Odds float64 `gorm:"not null" json:"odds"`
}

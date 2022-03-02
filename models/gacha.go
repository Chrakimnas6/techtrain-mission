package models

type Gacha struct {
	ID    uint `json:"id"`
	Times uint `gorm:"not null" json:"times"`
}

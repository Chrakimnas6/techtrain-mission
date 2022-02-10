package models

type Gacha struct {
	ID    uint `gorm:"autoIncrement"`
	Times uint `gorm:"not null" json:"times"`
}

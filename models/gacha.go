package models

type Gacha struct {
	ID    uint
	Times uint `gorm:"not null" json:"times"`
}

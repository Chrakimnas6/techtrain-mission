package models

type CharacterSSR struct {
	ID          uint `gorm:"autoIncrement"`
	CharacterID uint
}

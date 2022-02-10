package models

type UserCharacter struct {
	ID          uint `gorm:"autoIncrement"`
	UserID      uint
	CharacterID uint
}

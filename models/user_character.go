package models

type UserCharacter struct {
	ID            uint `gorm:"autoIncrement"`
	UserID        uint
	CharacterID   uint
	CharacterRank string
	Name          string
	Level         uint `gorm:"default:1"`
}

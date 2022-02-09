package models

type CharacterOdds struct {
	CharacterID uint `gorm:"primaryKey"`
	Odds        float64
}

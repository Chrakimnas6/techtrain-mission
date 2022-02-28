package models

type GachaCharacterOdds struct {
	GachaID     uint    `json:"gachaId"`
	CharacterID uint    `gorm:"foreignKey:ID"`
	Odds        float64 `gorm:"not null" json:"odds"`
}

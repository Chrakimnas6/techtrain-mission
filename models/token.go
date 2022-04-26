package models

type Token struct {
	ID      uint   `gorm:"autoIncrement"`
	Symbol  string `gorm:"not null" json:"symbol"`
	Address string
}

package models

type Transaction struct {
	ID     uint   `gorm:"autoIncrement"`
	Tx     string `gorm:"not null" json:"tx"`
	Status string `gorm:"not null" json:"status"`
}

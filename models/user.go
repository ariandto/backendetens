package models

type User struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	NIK        string `json:"nik" gorm:"unique;not null"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Shift      string `json:"shift"`
	Photo      string `json:"photo"` // nama file / path
	Phone      string `json:"phone"`
}

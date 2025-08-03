package models

type Product struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `json:"name"`
	Size     string  `json:"size"`
	Stock    int     `json:"stock"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
}

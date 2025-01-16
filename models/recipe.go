package models

type Recipe struct {
	ID           uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string   `gorm:"type:text;not null" json:"name"`
	Ingredients  []string `gorm:"serializer:json" json:"ingredients"`
	Instructions []string `gorm:"serializer:json" json:"instructions"`
	ImageUrl     string   `gorm:"type:text;not null" json:"imageUrl"`
}

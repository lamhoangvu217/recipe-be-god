package services

import (
	"recipe-be-god/database"
	"recipe-be-god/models"
)

func GetAllRecipes() ([]models.Recipe, error) {
	var recipes []models.Recipe
	if err := database.DB.Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func GetRecipesByCuisineId(cuisineId uint) ([]models.Recipe, error) {
	var recipes []models.Recipe
	query := database.DB.Where("cuisine_id = ?", cuisineId).Preload("Cuisine")
	if err := query.Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

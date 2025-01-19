package services

import (
	"errors"
	"gorm.io/gorm"
	"recipe-be-god/database"
	"recipe-be-god/models"
	"recipe-be-god/utils"
	"strings"
)

func GetAllRecipesService() ([]models.Recipe, error) {
	var recipes []models.Recipe
	if err := database.DB.Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func GetRecipesByCuisineIdService(cuisineId uint) ([]models.Recipe, error) {
	var recipes []models.Recipe
	query := database.DB.Where("cuisine_id = ?", cuisineId).Preload("Cuisine")
	if err := query.Find(&recipes).Error; err != nil {
		return nil, err
	}
	return recipes, nil
}

func CreateRecipeService(recipe *models.Recipe) (*models.Recipe, error) {
	if err := database.DB.Create(&recipe).Error; err != nil {
		return nil, err
	}
	return recipe, nil
}

func DeleteRecipeService(recipeId uint) error {
	// Start a transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Retrieve the recipe to ensure it exists
	var recipe models.Recipe
	if err := tx.First(&recipe, recipeId).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("recipe id not found")
		}
		return errors.New("could not retrieve recipe")
	}

	// Delete the recipe itself
	if err := tx.Delete(&recipe).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete recipe")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction commit failed")
	}

	return nil
}

func UpdateRecipeService(updatedRecipe *models.Recipe) error {
	if err := database.DB.Save(&updatedRecipe).Error; err != nil {
		return err
	}
	return nil
}

func GetRecipeByIdService(recipeId uint) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := database.DB.First(&recipe, recipeId).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

const (
	MaxLevenshteinDistance = 3 // Adjust this threshold as needed
)

func SearchRecipesService(query string) ([]models.Recipe, error) {
	var allRecipes []models.Recipe
	if err := database.DB.Find(&allRecipes).Error; err != nil {
		return nil, err
	}

	var matchedRecipes []models.Recipe
	normalizedQuery := utils.RemoveVietnameseTones(strings.TrimSpace(query))

	for _, recipe := range allRecipes {
		// Try exact matching first
		normalizedName := utils.RemoveVietnameseTones(recipe.Name)
		if strings.Contains(normalizedName, normalizedQuery) {
			matchedRecipes = append(matchedRecipes, recipe)
			continue
		}

		// Try fuzzy matching if exact match fails
		if utils.LevenshteinDistance(normalizedName, normalizedQuery) <= MaxLevenshteinDistance {
			matchedRecipes = append(matchedRecipes, recipe)
			continue
		}

		// Check ingredients
		for _, ingredient := range recipe.Ingredients {
			normalizedIngredient := utils.RemoveVietnameseTones(ingredient)
			if strings.Contains(normalizedIngredient, normalizedQuery) ||
				utils.LevenshteinDistance(normalizedIngredient, normalizedQuery) <= MaxLevenshteinDistance {
				matchedRecipes = append(matchedRecipes, recipe)
				break
			}
		}
	}

	return matchedRecipes, nil
}

func SearchRecipesByCuisineServices(cuisineId uint, query string) ([]models.Recipe, error) {
	// First, get recipes by cuisine ID
	var recipesInCuisine []models.Recipe
	if err := database.DB.Where("cuisine_id = ?", cuisineId).Preload("Cuisine").Find(&recipesInCuisine).Error; err != nil {
		return nil, err
	}

	// Then apply search filtering
	var matchedRecipes []models.Recipe
	normalizedQuery := utils.RemoveVietnameseTones(strings.TrimSpace(query))

	for _, recipe := range recipesInCuisine {
		// Try exact matching first
		normalizedName := utils.RemoveVietnameseTones(recipe.Name)
		if strings.Contains(normalizedName, normalizedQuery) {
			matchedRecipes = append(matchedRecipes, recipe)
			continue
		}

		// Try fuzzy matching if exact match fails
		if utils.LevenshteinDistance(normalizedName, normalizedQuery) <= MaxLevenshteinDistance {
			matchedRecipes = append(matchedRecipes, recipe)
			continue
		}

		// Check ingredients
		for _, ingredient := range recipe.Ingredients {
			normalizedIngredient := utils.RemoveVietnameseTones(ingredient)
			if strings.Contains(normalizedIngredient, normalizedQuery) ||
				utils.LevenshteinDistance(normalizedIngredient, normalizedQuery) <= MaxLevenshteinDistance {
				matchedRecipes = append(matchedRecipes, recipe)
				break
			}
		}
	}

	return matchedRecipes, nil
}

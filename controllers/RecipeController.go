package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
	"recipe-be-god/database"
	"recipe-be-god/models"
	"recipe-be-god/services"
	"strconv"
	"strings"
)

func GetRecipesByCuisine(c *fiber.Ctx) error {
	cuisineIdStr := c.Query("cuisineId")
	searchQuery := strings.TrimSpace(c.Query("search"))

	var recipes []models.Recipe
	var err error

	// If we have both cuisine ID and search query
	if cuisineIdStr != "" && searchQuery != "" {
		cuisineId, err := strconv.ParseUint(cuisineIdStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid cuisine id",
			})
		}
		recipes, err = services.SearchRecipesByCuisineServices(uint(cuisineId), searchQuery)
	} else if cuisineIdStr != "" {
		// Only cuisine filter
		cuisineId, err := strconv.ParseUint(cuisineIdStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid cuisine id",
			})
		}
		recipes, err = services.GetRecipesByCuisineIdService(uint(cuisineId))
	} else if searchQuery != "" {
		// Only search query
		recipes, err = services.SearchRecipesService(searchQuery)
	} else {
		// No filters - return all recipes
		recipes, err = services.GetAllRecipesService()
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "get recipes successfully",
		"recipes": recipes,
	})
}

func CreateRecipe(c *fiber.Ctx) error {
	recipe := new(models.Recipe)
	if err := c.BodyParser(recipe); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	createdRecipe, err := services.CreateRecipeService(recipe)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "recipe created successfully",
		"recipe":  createdRecipe,
	})
}

func DeleteRecipe(c *fiber.Ctx) error {
	recipeIdStr := c.Params("id")
	if recipeIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "recipe id is required",
		})
	}
	// Convert task id from string to uint
	recipeId, err := strconv.ParseUint(recipeIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid task id",
		})
	}
	if err := services.DeleteRecipeService(uint(recipeId)); err != nil {
		if err.Error() == "recipe id not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "recipe deleted successfully",
	})
}

func UpdateRecipe(c *fiber.Ctx) error {
	recipeIdStr := c.Params("id")
	if recipeIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "recipe id is required",
		})
	}
	// Convert product id from string to uint
	recipeId, err := strconv.ParseUint(recipeIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid task id",
		})
	}
	var recipe models.Recipe
	if err := database.DB.First(&recipe, recipeId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "recipe id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve recipe",
		})
	}
	var updateRecipeData models.Recipe
	if err := c.BodyParser(&updateRecipeData); err != nil {
		// Return 400 if request body is invalid
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if updateRecipeData.Name != "" {
		recipe.Name = updateRecipeData.Name
	}
	if updateRecipeData.Ingredients != nil {
		recipe.Ingredients = updateRecipeData.Ingredients
	}
	if updateRecipeData.Instructions != nil {
		recipe.Instructions = updateRecipeData.Instructions
	}
	if updateRecipeData.ImageUrl != "" {
		recipe.ImageUrl = updateRecipeData.ImageUrl
	}
	if updateRecipeData.CuisineID != 0 {
		recipe.CuisineID = updateRecipeData.CuisineID
	}
	if err := services.UpdateRecipeService(&recipe); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Recipe updated successfully",
		"recipe":  recipe,
	})
}

func GetRecipeById(c *fiber.Ctx) error {
	recipeIdStr := c.Params("id")
	if recipeIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "recipe id is required",
		})
	}
	// Convert product id from string to uint
	recipeId, err := strconv.ParseUint(recipeIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid recipe id",
		})
	}
	recipe, err := services.GetRecipeByIdService(uint(recipeId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "get recipe successfully",
		"recipe":  recipe,
	})
}

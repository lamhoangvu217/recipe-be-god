package controllers

import (
	"github.com/gofiber/fiber/v2"
	"recipe-be-god/services"
	"strconv"
)

func GetAllRecipes(c *fiber.Ctx) error {
	recipes, err := services.GetAllRecipes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "get all recipes successfully",
		"recipes": recipes,
	})
}

func GetRecipesByCuisine(c *fiber.Ctx) error {
	cuisineIdStr := c.Query("cuisineId")
	if cuisineIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cuisine id is required",
		})
	}
	// Convert categoryId from string to uint
	cuisineId, err := strconv.ParseUint(cuisineIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid cuisine id",
		})
	}
	recipes, err := services.GetRecipesByCuisineId(uint(cuisineId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":   "get recipes successfully",
		"recipes":   recipes,
		"cuisineId": cuisineId,
	})
}

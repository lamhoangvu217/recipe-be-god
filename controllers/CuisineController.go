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
)

func GetAllCuisine(c *fiber.Ctx) error {
	cuisines, err := services.GetAllCuisineService()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":  "get all cuisines successfully",
		"cuisines": cuisines,
	})
}

func CreateCuisine(c *fiber.Ctx) error {
	cuisine := new(models.Cuisine)
	if err := c.BodyParser(cuisine); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}
	createdCuisine, err := services.CreateCuisineService(cuisine)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "cuisine created successfully",
		"cuisine": createdCuisine,
	})
}

func DeleteCuisine(c *fiber.Ctx) error {
	cuisineIdStr := c.Params("id")
	if cuisineIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cuisine id is required",
		})
	}
	// Convert task id from string to uint
	cuisineId, err := strconv.ParseUint(cuisineIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid task id",
		})
	}
	if err := services.DeleteCuisineService(uint(cuisineId)); err != nil {
		if err.Error() == "cuisine id not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "cuisine and all associated recipes deleted successfully",
	})
}

func UpdateCuisine(c *fiber.Ctx) error {
	cuisineIdStr := c.Params("id")
	if cuisineIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cuisine id is required",
		})
	}
	// Convert product id from string to uint
	cuisineId, err := strconv.ParseUint(cuisineIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid cuisine id",
		})
	}
	var cuisine models.Cuisine
	if err := database.DB.First(&cuisine, cuisineId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "cuisine id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve cuisine",
		})
	}
	var updateCuisineData models.Cuisine
	if err := c.BodyParser(&updateCuisineData); err != nil {
		// Return 400 if request body is invalid
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if updateCuisineData.Name != "" {
		cuisine.Name = updateCuisineData.Name
	}

	if err := services.UpdateCuisineService(&cuisine); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cuisine updated successfully",
		"cuisine": cuisine,
	})
}

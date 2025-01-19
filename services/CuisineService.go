package services

import (
	"errors"
	"gorm.io/gorm"
	"recipe-be-god/database"
	"recipe-be-god/models"
)

func GetAllCuisineService() ([]models.Cuisine, error) {
	var cuisines []models.Cuisine
	if err := database.DB.Find(&cuisines).Error; err != nil {
		return nil, err
	}
	return cuisines, nil
}

func CreateCuisineService(cuisine *models.Cuisine) (*models.Cuisine, error) {
	if err := database.DB.Create(&cuisine).Error; err != nil {
		return nil, err
	}
	return cuisine, nil
}

func DeleteCuisineService(cuisineId uint) error {
	// Start a transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Retrieve the cuisine to ensure it exists
	var cuisine models.Cuisine
	if err := tx.First(&cuisine, cuisineId).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("cuisine id not found")
		}
		return errors.New("could not retrieve cuisine")
	}

	// Delete all recipe associated with the cuisine
	if err := tx.Where("cuisine_id = ?", cuisineId).Delete(&models.Recipe{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete recipe")
	}

	// Delete the cuisine itself
	if err := tx.Delete(&cuisine).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete cuisine")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.New("transaction commit failed")
	}

	return nil
}

func UpdateCuisineService(updatedCuisine *models.Cuisine) error {
	if err := database.DB.Save(&updatedCuisine).Error; err != nil {
		return err
	}
	return nil
}

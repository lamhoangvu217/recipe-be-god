package routes

import (
	"github.com/gofiber/fiber/v2"
	"recipe-be-god/controllers"
)

func Setup(app *fiber.App) {
	app.Get("/api/recipes", controllers.GetRecipesByCuisine)
}

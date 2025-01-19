package routes

import (
	"github.com/gofiber/fiber/v2"
	"recipe-be-god/controllers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/recipes", controllers.GetRecipesByCuisine)
	api.Get("/cuisines", controllers.GetAllCuisine)

	api.Post("/cuisine", controllers.CreateCuisine)
	api.Delete("/cuisine/:id", controllers.DeleteCuisine)
	api.Put("/cuisine/:id", controllers.UpdateCuisine)

	api.Post("/recipe", controllers.CreateRecipe)
	api.Get("/recipe/:id", controllers.GetRecipeById)
	api.Delete("/recipe/:id", controllers.DeleteRecipe)
	api.Put("/recipe/:id", controllers.UpdateRecipe)
}

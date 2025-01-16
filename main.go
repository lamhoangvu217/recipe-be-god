package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"recipe-be-god/database"
	"recipe-be-god/routes"
)

func main() {
	database.Connect()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3000", // Replace with the actual frontend origin(s)
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	routes.Setup(app)
	app.Listen(":" + port)

}

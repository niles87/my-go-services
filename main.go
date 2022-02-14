package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

type Message struct {
	Msg string
}

type User struct {
	Id       int64
	Name     string
	Email    string
	Password string
	Wins     int
	Losses   int
	Draws    int
}

func main() {
	// Load dotenv file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Print("ENV failed to load")
	}

	connect()

	// Set port (for heroku later)
	PORT := os.Getenv("PORT")

	// use Golangs default html templates
	engine := html.New("./views", ".html")

	// Initialize app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")

	// Add middleware with .Use
	app.Use(logger.New())
	app.Use(requestid.New())
	routes(app)

	log.Fatal(app.Listen(":" + PORT))
}

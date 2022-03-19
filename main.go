package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"github.com/niles87/my-go-services/controller"
	"github.com/niles87/my-go-services/mydb"
)

func main() {
	// Load dotenv file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Print("ENV failed to load")
	}

	db, err := mydb.Connect()
	if err != nil {
		log.Fatal(err)
	}

	hdl := controller.NewDBHandler(db)

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
	// Display index route
	app.Get("/", hdl.GetIndex)

	// Group related endpoints together
	userApp := app.Group("/user")
	userApp.Get("", hdl.GetUsers)
	userApp.Post("", hdl.CreateUser)
	userApp.Put("", hdl.UpdateUser)
	userApp.Get("/:id", hdl.GetUserById)
	userApp.Delete("/:id", hdl.DeleteUser)
	userApp.Post("/login", hdl.Login)

	log.Fatal(app.Listen(":" + PORT))
}

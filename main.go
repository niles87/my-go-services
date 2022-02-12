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

func getIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":  "Welcome",
		"Second": "To my services",
	})
}

type User struct {
	Name  string
	Email string
}

func getUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(user)
}

var user User

func createUser(c *fiber.Ctx) error {

	body := new(User)

	err := c.BodyParser(body)

	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(Message{Msg: err.Error()})
		return err
	}

	user = User{
		Name:  body.Name,
		Email: body.Email,
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func updateUser(c *fiber.Ctx) error {

	body := new(User)
	err := c.BodyParser(body)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(Message{Msg: err.Error()})
		return err
	}

	user = *body
	return c.Status(fiber.StatusAccepted).JSON(user)
}

func main() {
	// Load dotenv file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Print("ENV failed to load")
	}

	// Set port (for heroku later)
	PORT := os.Getenv("PORT")

	// use Golangs default html templates
	engine := html.New("./views", ".html")

	// Initialize app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")

	// Display index route
	app.Get("/", getIndex)

	// Add middleware with .Use
	app.Use(logger.New())
	app.Use(requestid.New())

	// Group related endpoints together
	userApp := app.Group("/user")
	userApp.Get("", getUser)
	userApp.Post("", createUser)
	userApp.Put("", updateUser)

	log.Fatal(app.Listen(":" + PORT))
}

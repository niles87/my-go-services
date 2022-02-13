package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":  "Welcome",
		"Second": "To my services",
	})
}

func getUsers(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(users)
}

func getUser(c *fiber.Ctx) error {
	i, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(Message{Msg: "Missing params"})
	}
	for _, val := range users {
		if val.Id == int64(i) {
			return c.Status(fiber.StatusOK).JSON(val)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(Message{Msg: "User not found"})
}

func createUser(c *fiber.Ctx) error {

	body := new(User)

	err := c.BodyParser(body)

	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(Message{Msg: err.Error()})
		return err
	}

	user := User{
		Id:       int64(len(users) + 1),
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
		Wins:     0,
		Losses:   0,
		Draws:    0,
	}

	users = append(users, user)

	return c.Status(fiber.StatusCreated).JSON(user)
}

func updateUser(c *fiber.Ctx) error {

	body := new(User)
	err := c.BodyParser(body)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(Message{Msg: err.Error()})
		return err
	}

	user := *body

	for i, val := range users {
		if int64(val.Id) == user.Id {
			users[i] = user
			return c.Status(fiber.StatusAccepted).JSON(user)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(Message{Msg: "User Not Found"})
}

func routes(app *fiber.App) {
	// Display index route
	app.Get("/", getIndex)

	// Group related endpoints together
	userApp := app.Group("/user")
	userApp.Get("", getUsers)
	userApp.Post("", createUser)
	userApp.Put("", updateUser)
	userApp.Get("/:id", getUser)

}
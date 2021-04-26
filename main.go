package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ziemedee/gofiber-learn/database"
	"github.com/ziemedee/gofiber-learn/routes"
)

func setUpRoutes(app *fiber.App) {
	app.Post("/register", routes.Register)
	app.Post("/login", routes.Login)
	app.Get("/employees", routes.GetEmployees)
	app.Get("/employee/:id", routes.GetEmployee)
	app.Post("/employee", routes.AddEmployee)
	app.Put("/employee/:id", routes.UpdateEmployees)
	app.Delete("/employee/:id", routes.DeleteEmployees)
}

func main() {

	if err := database.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	setUpRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("go fiber")
	})

	log.Fatal(app.Listen(":3000"))

}

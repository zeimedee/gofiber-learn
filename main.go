package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ziemedee/gofiber-learn/database"
	"github.com/ziemedee/gofiber-learn/middleware"
	"github.com/ziemedee/gofiber-learn/routes"
)

func setUpRoutes(app *fiber.App) {
	app.Post("/register", middleware.Auth(), routes.Register)
	app.Post("/login", routes.Login)
	app.Get("/employees", middleware.Auth(), routes.GetEmployees)
	app.Get("/employee/:id", middleware.Auth(), routes.GetEmployee)
	app.Post("/employee", middleware.Auth(), routes.AddEmployee)
	app.Put("/employee/:id", middleware.Auth(), routes.UpdateEmployees)
	app.Delete("/employee/:id", middleware.Auth(), routes.DeleteEmployees)
}

func main() {

	if err := database.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE",
		AllowHeaders:     "",
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           999999,
	}))
	setUpRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("go fiber")
	})

	log.Fatal(app.Listen(":4000"))

}

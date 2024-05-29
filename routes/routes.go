package routes

import (
	c "cafeshop-backend/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func User(app *fiber.App) {
	app.Use(logger.New())
	api := app.Group("/api")
	// /v1
	v1 := api.Group("/user")
	v1.Post("/login", c.Login)
	v1.Post("/register", c.Register)
	v1.Post("/addOrder/:id", c.CreateOrder)

}
func Admin(app *fiber.App) {
	app.Use(logger.New())
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "1234"},
	}))
	api := app.Group("/api")
	// /v2
	v2 := api.Group("/admin")
	v2.Get("/user", c.GetUser)
	v2.Put("/user/:id", c.UpdateUser)
	v2.Delete("/user/:id", c.DeleteUser)

}
func Product(app *fiber.App) {
	app.Use(logger.New())
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "1234",
		},
	}))
	api := app.Group("/api")
	// /v3
	v3 := api.Group("/Product")
	v3.Get("/", c.GetProduct)
	v3.Post("/", c.CreateProduct)
	v3.Put("/:id", c.UpdatetProduct)
	v3.Delete("/:id", c.DeletetProduct)
}
func Order(app *fiber.App) {
	app.Use(logger.New())
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "1234",
		},
	}))
	api := app.Group("/api")
	// /v3
	v4 := api.Group("/Product")
	v4.Get("/getProduct", c.GetProductsAndOrders)
}

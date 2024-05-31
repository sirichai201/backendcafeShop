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
	v1.Put("/updateUser/:id", c.UpdateUser)
}
func Product(app *fiber.App) {
	app.Use(logger.New())

	api := app.Group("/api")

	v5 := api.Group("/Product")
	v5.Get("/Product", c.GetProduct)
	v5.Get("/GetProductByName", c.GetProductByName)

	// v5.Get("/test", c.GetFacID(database.ProductsCollection))

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
	v2.Get("/GetProduct", c.GetProduct)
	v2.Get("/GetProductByID/:id", c.GetProductByID)

	v2.Post("/", c.CreateProduct)
	v2.Put("/:id", c.UpdatetProduct)
	v2.Delete("/:id", c.DeletetProduct)
	v2.Get("/user", c.GetUser)
	v2.Put("/user/:id", c.UpdateUser)
	v2.Delete("/user/:id", c.DeleteUser)

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
	v4 := api.Group("/OrderProduct")
	v4.Get("/GetProductsAndOrders", c.GetProductsAndOrders)
	v4.Get("/GetProductOrders/:id", c.GetProductOrders)
	v4.Put("/ChackOut/:id", c.ChackOut)
	v4.Post("/CreateBill/:id", c.CreateBillByUserID)
}

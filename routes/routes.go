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

	//user edit
	v1 := api.Group("/user")
	v1.Post("/login", c.Login)
	v1.Post("/register", c.Register)
	v1.Post("/addOrder/:id", c.AddOrder)
	v1.Put("/updateUser/:id", c.UpdateUser)
	v1.Put("/ChackOut/:id", c.CheckOut)
	v1.Post("/CreateBill/:id", c.CreateBillByUserID)
}
func Product(app *fiber.App) {
	app.Use(logger.New())

	api := app.Group("/api")


	// Get Product
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

	//admin edit product
	v2.Get("/GetProduct", c.GetProduct)
	v2.Get("/GetProductByID/:id", c.GetProductByID)
	v2.Post("/PostProduct", c.CreateProduct)
	v2.Put("/UpdateProduct/:id", c.UpdatetProduct)
	v2.Delete("/DeleteProduct/:id", c.DeletetProduct)

   //admin edit user 
	v2.Get("/user", c.GetUser)
	v2.Put("/UpdateUser/:id", c.UpdateUser)
	v2.Delete("/DeleteUser/:id", c.DeleteUser)

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
	v4.Get("/GetBill",c.GetBill)

}

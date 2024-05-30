// In main.go

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"cafeshop-backend/database"
	"cafeshop-backend/routes"
)

const mongoURI = "mongodb+srv://sirichaichantharasri4:First0903319646@cafeshop.iy0znlw.mongodb.net/?retryWrites=true&w=majority&appName=cafeShop"

func main() {
	// Connect to MongoDB
	if err := database.InitDatabase(mongoURI); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	// Initialize collections
	database.InitCollection()

	// Initialize the Fiber app
	app := fiber.New()


	// Use the logger middleware
	app.Use(logger.New())

	// Register routes
	routes.User(app)
	routes.Product(app)
	routes.Admin(app)
	routes.Order(app)

	// Start the Fiber server
	log.Fatal(app.Listen(":3000"))
}

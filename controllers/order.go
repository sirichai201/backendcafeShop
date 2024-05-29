package controllers

import (
	"cafeshop-backend/database"
	"cafeshop-backend/models"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	order := new(models.Orders)
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	product := new(models.Products)
	err = database.ProductsCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).SendString("Product not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	price := product.ProductPrice

	order.Total_Price = price * float64(order.Quantity)

	tel := order.Tel
	var user models.Users
	err = database.UsersCollection.FindOne(context.Background(), bson.M{"phone": tel}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	order.User_id = user.User_ID
	order.Product_id = objID

	result, err := database.OrdersCollection.InsertOne(context.Background(), order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	order.Order_id = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(order)
}

func GetProductsAndOrders(c *fiber.Ctx) error {
	// Find all products
	cursor, err := database.ProductsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer cursor.Close(context.Background())

	// Variables to store product with orders and count of products with orders
	var products []models.Products
	productsWithOrdersCount := 0

	for cursor.Next(context.Background()) {
		var product models.Products
		if err := cursor.Decode(&product); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Check if there are orders for the product
		ordersCursor, err := database.OrdersCollection.Find(context.Background(), bson.M{"product_id": product.Product_ID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		defer ordersCursor.Close(context.Background())

		if ordersCursor.Next(context.Background()) {
			products = append(products, product)
			productsWithOrdersCount++
		}
	}

	// Prepare the response
	type ProductsSummary struct {
		ProductsWithOrdersCount int               `json:"products_with_orders_count"`
		Products                []models.Products `json:"products"`
	}
	summary := ProductsSummary{
		ProductsWithOrdersCount: productsWithOrdersCount,
		Products:                products,
	}

	return c.Status(fiber.StatusOK).JSON(summary)
}

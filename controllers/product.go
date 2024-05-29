package controllers

import (
	"cafeshop-backend/database"
	"cafeshop-backend/models"
	"context"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProduct(c *fiber.Ctx) error {
	cursor, err := database.ProductsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer cursor.Close(context.Background())

	var Product []models.Products
	if err := cursor.All(context.Background(), &Product); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(Product)
}
func CreateProduct(c *fiber.Ctx) error {
	Product := new(models.Products)

	if err := c.BodyParser(Product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	validate := validator.New()
	errors := validate.Struct(Product)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}
	var existingProduct models.Products
	err := database.UsersCollection.FindOne(context.Background(), bson.M{"product_name": Product.ProductName}).Decode(&existingProduct)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Product_Name already exists")
	}

	result, err := database.ProductsCollection.InsertOne(context.Background(), Product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	Product.Product_ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(Product)
}

func UpdatetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Product ID")
	}

	Product := new(models.Products)
	if err := c.BodyParser(Product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err = database.ProductsCollection.ReplaceOne(context.Background(), bson.M{"_id": objID}, Product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
func DeletetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Product ID")
	}

	_, err = database.ProductsCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

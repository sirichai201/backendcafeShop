package controllers

import (
	"cafeshop-backend/database"
	"cafeshop-backend/models"
	"context"
	"log"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"

	// "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// func GetProduct(c *fiber.Ctx) error {
// 	cursor, err := database.ProductsCollection.Find(context.Background(), bson.M{})
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
// 	}
// 	defer cursor.Close(context.Background())

// 	var Product []models.Products
// 	if err := cursor.All(context.Background(), &Product); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
// 	}
// 	return c.JSON(Product)
// }

func GetProductByName(c *fiber.Ctx) error {
	name := c.Query("name")

	if name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Query parameter 'name' is required")
	}

	var product models.Products
	log.Printf("Searching for product with name: %s\n", name)
	filter := bson.M{"productname": name}
	log.Printf("Filter used: %v\n", filter)

	err := database.ProductsCollection.FindOne(context.Background(), filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Product not found:", err)
			return c.Status(fiber.StatusNotFound).SendString("Product not found")
		}
		log.Println("Error occurred while searching for product:", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	log.Printf("Product found: %+v\n", product)
	return c.JSON(product)
}

func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid Product ID:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Product ID")
	}

	var product models.Products
	err = database.ProductsCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		log.Println("Product not found:", err)
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	return c.JSON(product)
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
	err := database.ProductsCollection.FindOne(context.Background(), bson.M{"productname": Product.ProductName}).Decode(&existingProduct)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Product_Name already exists")
	}

	if Product.ProductPrice < 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Product_Price must be greater than 0")
	}

	_, err = database.ProductsCollection.InsertOne(context.Background(), Product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(Product)
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Products

	cursor, err := database.ProductsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &products); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func UpdatetProduct(c *fiber.Ctx) error {
	productID := c.Params("id")
	var updateData models.Products

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	validate := validator.New()
	errors := validate.Struct(updateData)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}

	if updateData.ProductPrice < 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Product_Price must be greater than 0")
	}

	update := bson.M{
		"$set": bson.M{
			"productname":  updateData.ProductName,
			"productprice": updateData.ProductPrice,
			"productpoint": updateData.ProductPoint,
			"producttype":  updateData.ProductType,
			"description":  updateData.Description,
			"image":        updateData.Image,
		},
	}

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	_, err = database.ProductsCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(updateData)
}

func DeletetProduct(c *fiber.Ctx) error {
	productID := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	_, err = database.ProductsCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

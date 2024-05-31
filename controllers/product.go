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

// func GetProductByName(c *fiber.Ctx) error {
// 	name := c.Params("name")
// 	var product models.Products
// 	log.Printf("Searching for product with name: %s\n", name)
// 	productname := bson.M{"product_name": name}
// 	log.Printf("Filter used: %v\n", productname)

//		database.ProductsCollection.FindOne(context.Background(), bson.M{"product_name": productname}).Decode(&product)
//		log.Printf("Product found: %+v\n", product)
//		return c.JSON(product)
//	}
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
	err := database.UsersCollection.FindOne(context.Background(), bson.M{"productname": Product.ProductName}).Decode(&existingProduct)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Product_Name already exists")
	}
	if Product.ProductPrice < 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Price must more than 0")
	}
	if Product.ProductPoint < 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Point must more than 0")
	}
	result, err := database.ProductsCollection.InsertOne(context.Background(), Product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	Product.ProductID = result.InsertedID.(primitive.ObjectID)

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
	if Product.ProductPrice < 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Price must more than 0")
	}
	if Product.ProductPoint < 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Point must more than 0")
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

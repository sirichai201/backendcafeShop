package controllers

import (
	"cafeshop-backend/database"
	"cafeshop-backend/models"
	"context"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var store = session.New()

func Login(c *fiber.Ctx) error {
	var loginData struct {
		UserName string `json:"username" validate:"required,min=3,max=32"`
		Password string `json:"password" validate:"required,min=3,max=32"`
	}
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	validate := validator.New()
	if err := validate.Struct(loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	user := new(models.Users)
	err := database.UsersCollection.FindOne(context.Background(), bson.M{"username": loginData.UserName}).Decode(user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid username or password")
	}
	if user.Status != 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Account Blocked")
	}

	// Create a response struct to include the user's name and phone number
	response := struct {
		Message     string `json:"message"`
		UserName    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}{
		Message:     "Login Success",
		UserName:    user.UserName, // Assuming the field name is UserName in models.Users
		PhoneNumber: user.Phone,
		Email:       user.Email, // Assuming the field name is PhoneNumber in models.Users
	}

	return c.Status(fiber.StatusAccepted).JSON(response)
}

// func Login(c *fiber.Ctx) error {
// 	var loginData struct {
// 		UserName string `json:"username" validate:"required,min=3,max=32"`
// 		Password string `json:"password" validate:"required,min=3,max=32"`
// 	}
// 	if err := c.BodyParser(&loginData); err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}
// 	validate := validator.New()
// 	if err := validate.Struct(loginData); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
// 	}
// 	user := new(models.Users)
// 	err := database.UsersCollection.FindOne(context.Background(), bson.M{"username": loginData.UserName}).Decode(user)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).SendString("Invalid username or password")
// 	}
// 	if user.Status != 1 {
// 		return c.Status(fiber.StatusBadRequest).SendString("Account Blocked")
// 	}

// 	sess, err := store.Get(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create session")
// 	}

// 	// บันทึก User ID ใน session
// 	sess.Set("user_id", user.User_ID.Hex())

// 	// ตรวจสอบและบันทึก session
// 	if err := sess.Save(); err != nil {
// 		fmt.Printf("Error saving session: %v\n", err) // พิมพ์ข้อผิดพลาดลงในคอนโซล
// 		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
// 	}
// 	// Return session information in response
// 	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
// 		"message": "Login Success",
// 		"session": sess.Get("user_id"),
// 		"user_id": user.User_ID,
// 	})
// }

func Register(c *fiber.Ctx) error {
	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	user.Role = "User"
	user.Status = 1
	var existingUser models.Users
	err := database.UsersCollection.FindOne(context.Background(), bson.M{"username": user.UserName}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Username already exists")
	}
	var existingUserByEmail models.Users
	err = database.UsersCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUserByEmail)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Email already exists")
	}
	var existingUserByPhone models.Users
	err = database.UsersCollection.FindOne(context.Background(), bson.M{"phone": user.Phone}).Decode(&existingUserByPhone)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).SendString("Phone already exists")
	}

	validate := validator.New()
	errors := validate.Struct(user)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}

	result, err := database.UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	user.User_ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(user)
}

package controllers

import (
	"cafeshop-backend/database"
	"cafeshop-backend/models"
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("รหัสผลิตภัณฑ์ไม่ถูกต้อง")
	}

	order := new(models.Orders)
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if order.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("จำนวนต้องมากกว่า 0")
	}

	if order.Member == "" {
		return c.Status(fiber.StatusBadRequest).SendString("กรุณาใส่ข้อมูลสมาชิก")
	}

	product := new(models.Products)

	// โหลดตำแหน่งเวลา Bangkok
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("ไม่สามารถโหลดตำแหน่งเวลาได้")
	}
	// ตั้งค่าเวลาปัจจุบันตามเขตเวลา Bangkok
	order.CreatedAt = time.Now().In(location)

	err = database.ProductsCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).SendString("ไม่พบผลิตภัณฑ์")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	order.Total_Price = int(product.ProductPrice * float64(order.Quantity))

	member := order.Member
	var user models.Users
	err = database.UsersCollection.FindOne(context.Background(), bson.M{"phone": member}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).SendString("ไม่พบผู้ใช้")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	order.PesentPoint = int(float64(product.ProductPoint) * float64(order.Quantity))
	user.Point += int(float64(product.ProductPoint) * float64(order.Quantity))

	_, err = database.UsersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.User_ID},
		bson.M{"$set": bson.M{"point": user.Point}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	order.User_id = user.User_ID
	order.Product_id = objID
	order.Status = "0"

	// Insert the document into the database
	result, err := database.OrdersCollection.InsertOne(context.Background(), order)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	order.Order_id = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(order)
}
func GetProductsAndOrders(c *fiber.Ctx) error { //เช็คProduct  ที่มี  Order โชว์รายการ product ที่มี ออเดอร์ออกมา
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
		ordersCursor, err := database.OrdersCollection.Find(context.Background(), bson.M{"product_id": product.ProductID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		defer ordersCursor.Close(context.Background())

		// Check if there is at least one order for the product
		if ordersCursor.RemainingBatchLength() > 0 {
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

func GetProductOrders(c *fiber.Ctx) error { //เช็ค Order  ใน  Product  ว่าProductนั้นๆมี Order อะไรบ้าง
	// รับค่า product ID จากพารามิเตอร์
	productID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	// Find the product
	var product models.Products
	err = database.ProductsCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	// Find orders for the product
	var orders []models.Orders
	cursor, err := database.OrdersCollection.Find(context.Background(), bson.M{"product_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var order models.Orders
		if err := cursor.Decode(&order); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		orders = append(orders, order)
	}

	// Prepare the response
	type ProductOrderSummary struct {
		Product     models.Products `json:"product"`
		OrdersCount int             `json:"orders_count"`
		Orders      []models.Orders `json:"orders"`
	}

	summary := ProductOrderSummary{
		Product:     product,
		OrdersCount: len(orders),
		Orders:      orders,
	}

	return c.Status(fiber.StatusOK).JSON(summary)
}

func CreateBillByUserID(c *fiber.Ctx) error { //สร้างบิลที่ลูกค้าสั่งสินค้าจาก oder ค้นหา order จาก user_id และ status
	bill := new(models.Bill)
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid Product ID:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Product ID")
	}
	var orders []models.Orders
	filter := bson.M{
		"user_id": objID,
		"status":  "0", 
	}

	cursor, err := database.OrdersCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &orders); err != nil {
		log.Fatal(err)
	}
	product := new(models.Products)
	totalPrice := 0
	totalPoints := 0
	for _, order := range orders {
		product_id := order.Product_id
		err = database.ProductsCollection.FindOne(context.Background(), bson.M{"_id": product_id}).Decode(&product)
		log.Println("Price :", product.ProductPrice)
		log.Println("Points :", product.ProductPoint)
		totalPrice += int(product.ProductPrice) * order.Quantity
		totalPoints += product.ProductPoint * order.Quantity
	}

	// Create the bill
	bill = &models.Bill{
		User_id:     objID,
		CreatedAt:   time.Now(),
		Total_Price: totalPrice,
		PesentPoint: totalPoints,
		Payment:     "0", // Assuming you have a function to calculate the points
	}

	// Save the bill to the database
	result, err := database.BillsCollection.InsertOne(context.Background(), bill) // Assuming BillsCollection is the correct collection
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result :", result)

	return c.Status(fiber.StatusCreated).JSON(result)
}

func ChackOut(c *fiber.Ctx) error {
	id := c.Params("id") // Bill ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Bill ID")
	}

	// Find the bill with the specified ID and payment status "0"
	var bill models.Bill
	filter := bson.M{"_id": objID, "payment": "0"}
	err = database.BillsCollection.FindOne(context.Background(), filter).Decode(&bill)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).SendString("Bill not found")
		}
		log.Fatal(err)
	}

	// Update bill payment status to "1"
	bill.Payment = "1"
	_, err = database.BillsCollection.ReplaceOne(context.Background(), bson.M{"_id": objID}, bill)
	if err != nil {
		log.Fatal(err)
	}

	// Find all orders for the user associated with the bill
	var orders []models.Orders
	filter = bson.M{"user_id": bill.User_id}
	cursor, err := database.OrdersCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	if err = cursor.All(context.Background(), &orders); err != nil {
		log.Fatal(err)
	}

	// Update each order's status to "1"
	for _, order := range orders {
		log.Println(order)
		updateFilter := bson.M{"status": "0"}
		order.Status = "1"
		_, err = database.OrdersCollection.ReplaceOne(context.Background(), updateFilter, order)
		if err != nil {
			log.Fatal(err)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

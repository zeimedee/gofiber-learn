package routes

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/ziemedee/gofiber-learn/database"
	"github.com/ziemedee/gofiber-learn/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const Secret = "secret"

func Register(c *fiber.Ctx) error {
	collection := database.Mg.Db.Collection("admin")
	admin := new(models.Admin)

	if err := c.BodyParser(admin); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	query := bson.D{{Key: "user", Value: admin.User}}
	result := collection.FindOne(c.Context(), query)
	user := &models.Admin{}
	result.Decode(user)
	if user.User != "" {
		return c.Status(400).JSON("user exists")
	}
	insert, err := collection.InsertOne(c.Context(), admin)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	filter := bson.D{{Key: "_id", Value: insert.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdAdmin := &models.Admin{}
	createdRecord.Decode(createdAdmin)
	return c.Status(200).JSON(createdAdmin)
}

func Login(c *fiber.Ctx) error {
	collection := database.Mg.Db.Collection("admin")
	admin := new(models.Admin)
	if err := c.BodyParser(admin); err != nil {
		return c.Status(400).JSON("failed")
	}

	query := bson.D{{Key: "user", Value: admin.User}}
	result := collection.FindOne(c.Context(), query)

	ad := &models.Admin{}
	result.Decode(ad)
	if ad.User != "" && ad.Password != "" {

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		// claims["sub"] = "1"
		claims["exp"] = time.Now().Add(time.Hour * 24 * 7) // a week

		s, err := token.SignedString([]byte(Secret))
		if err != nil {
			return c.Status(500).JSON("internal server error")
		}

		cookie := new(fiber.Cookie)
		cookie.Name = "token"
		cookie.Value = s
		cookie.HTTPOnly = false
		cookie.MaxAge = 999999
		cookie.Expires = time.Now().Add(24 * time.Hour)
		// cookie.SameSite = "Lax"
		// cookie.Secure = true
		c.Cookie(cookie)

		return c.Status(200).JSON(fiber.Map{
			"user": ad.User,
		})
	}
	return c.Status(404).JSON("admin does not exist")
}

func GetEmployees(c *fiber.Ctx) error {
	query := bson.D{{}}

	cursor, err := database.Mg.Db.Collection("employees").Find(c.Context(), query)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var employees []models.Employee = make([]models.Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.Status(200).JSON(employees)
}

func GetEmployee(c *fiber.Ctx) error {
	employeeId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(404).SendString(err.Error())
	}
	query := bson.D{{Key: "_id", Value: employeeId}}
	result := database.Mg.Db.Collection("employees").FindOne(c.Context(), query)

	employee := &models.Employee{}
	result.Decode(employee)
	return c.Status(200).JSON(employee)
}

func AddEmployee(c *fiber.Ctx) error {
	collection := database.Mg.Db.Collection("employees")
	employee := new(models.Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	employee.Id = ""

	insertionResult, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdEmployee := &models.Employee{}
	createdRecord.Decode(createdEmployee)
	return c.Status(200).JSON(createdEmployee)
}

func UpdateEmployees(c *fiber.Ctx) error {
	collection := database.Mg.Db.Collection("employees")
	employeeId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.SendStatus(400)
	}
	employee := new(models.Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	query := bson.D{{Key: "_id", Value: employeeId}}
	update := bson.D{{
		Key: "$set", Value: bson.D{
			{Key: "name", Value: employee.Name},
			{Key: "salary", Value: employee.Salary},
			{Key: "age", Value: employee.Age},
		},
	}}

	err = collection.FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(400)
		}
		return c.SendStatus(500)
	}
	employee.Id = c.Params("id")
	return c.Status(200).JSON(employee)
}

func DeleteEmployees(c *fiber.Ctx) error {
	collection := database.Mg.Db.Collection("employees")
	employeeId, err := primitive.ObjectIDFromHex(c.Params("id"))

	query := bson.D{{Key: "_id", Value: employeeId}}
	result, err := collection.DeleteOne(c.Context(), query)
	if err != nil {
		return c.Status(401).SendString(err.Error())
	}

	if result.DeletedCount < 1 {
		return c.SendStatus(400)
	}

	return c.SendStatus(204)
}

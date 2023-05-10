package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "@Jemuel15"
	dbname   = "golang"
)

var DB *gorm.DB

func main() {
	migration()
	app := fiber.New()
	app.Get("/", GetUsers)
	app.Post("/", CreateUser)
	app.Listen(":3000")
}

func migration() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, _ = gorm.Open(postgres.Open(psqlconn), &gorm.Config{})
	DB.AutoMigrate(&User{})
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetUsers(c *fiber.Ctx) error {
	var users []User
	DB.Find(&users)
	return c.JSON(users)
}

func CreateUser(c *fiber.Ctx) error {
	// Create a new user record with the input values
	user := User{
		Name:    "John Doe",
		Email:   "johndoe@example.com",
		School:  "Example School",
		Company: "Example Company",
	}

	// Insert the new user record into the database
	result := DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": result.Error.Error()})
	}

	// Return the new user record as a JSON response
	return c.JSON(user)
}

type User struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	School  string `json:"school"`
	Company string `json:"company"`
}

type Tabless struct {
	Name string `json:"name"`
}

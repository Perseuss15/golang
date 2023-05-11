package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2" // Import Fiber package
	_ "github.com/lib/pq"         // Import Postgres driver
	"gorm.io/driver/postgres"     // Import Postgres GORM driver
	"gorm.io/gorm"                // Import GORM package
)

//	config := &Configuration{
//		host: os.Getenv("DB_HOST"),
//		port: os.Getenv("DB_PORT"),
//		user: os.Getenv("DB_USER"),
//		dbname: os.Getenv("DB_NAME"),
//		password: os.Getenv("DB_HOST")
//		sslmode: os.Getenv("DB_SSLMODE")
//	}
const (
	host     = "localhost"
	port     = 5435
	user     = "postgres"
	password = "postgres"
	dbname   = "golang"
)

var DB *gorm.DB // Declare a global variable to hold the database connection

func main() {
	migration()        // Call migration function to connect to the database
	app := fiber.New() // Create a new Fiber app instance
	api := app.Group("/api")
	api.Get("/getuser", GEtUseers)      // Define a route that responds to GET requests
	api.Post("/insertuser", CreateUser) // Define a route that responds to Post requests
	api.Post("/users/:id", DeleteUser)

	// api2 := app.Group("")
	app.Listen(":3000") // Start the server on port 3000
}

func migration() {
	// Define the database connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open a new database connection using GORM and the Postgres driver
	DB, _ = gorm.Open(postgres.Open(psqlconn), &gorm.Config{})
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GEtUseers(c *fiber.Ctx) error {
	// Declare a variable to hold the list of users
	var user []User

	// Use GORM to fetch all records from the "users" table
	DB.Find(&user)

	// Return the list of users as a JSON response
	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	// Parse the request body and extract the user data
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if the user already exists in the database
	existingUser := User{}
	result := DB.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		// User already exists, return an error response
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already exists"})
	}

	// Insert the user data into the database
	result = DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	// Return a success response
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	// Get the user ID from the request parameters
	userID := c.Params("id")
	fmt.Println("user ID", userID)

	// Delete the user from the database
	result := DB.Where("id = ?", userID).Delete(&User{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": result.Error.Error()})
	}

	// Return a success response
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

// Define a struct to represent a user record in the "users" table
type User struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	School  string `json:"school"`
	Company string `json:"company"`
}

// Define a struct to represent a table name
type Tabless struct {
	Name string `json:"name"`
}

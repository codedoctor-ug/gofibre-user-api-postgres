package main

import (
	"log"
	"os"

	"github.com/abiiranathan/gofibre-gorm-relations/db"
	"github.com/abiiranathan/gofibre-gorm-relations/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setUpRoutes(app *fiber.App) {
	app.Get("/api/users", models.GetUsers)
	app.Post("/api/users", models.CreateUser)
}

func ConnectDatabase() {
	dsn := os.Getenv("DSN")
	DBConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Set Database connection
	db.DBConn = DBConn

	if err != nil {
		log.Fatal("Unable to connect to the database: " + err.Error())
		return
	}

	dberr := DBConn.AutoMigrate(&models.User{}, &models.Profile{})

	if dberr != nil {
		panic("Error creating database migrations: " + dberr.Error())
	}
	log.Print("Connected to database")
}

func init() {
	godotenv.Load()
	ConnectDatabase()
}

func main() {
	app := fiber.New()
	setUpRoutes(app)
	log.Fatal(app.Listen(":3000"))
}

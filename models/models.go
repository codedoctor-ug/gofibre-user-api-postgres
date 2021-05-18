package models

import (
	"time"

	"github.com/abiiranathan/gofibre-gorm-relations/db"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// `validate:"required,email,min=6,max=32"`

type Profile struct {
	ID        int       `json:"id" gorm:"not null;primaryKey"`
	FullName  string    `json:"name" gorm:"size:50;not null" validate:"required"`
	Mobile    string    `json:"mobile" gorm:"size:10;not null;unique" validate:"required,len=10"`
	BirthDate time.Time `json:"birth_date" gorm:"not null" validate:"required"`
}

type User struct {
	ID       int     `json:"id" gorm:"not null;primaryKey"`
	Username string  `json:"username" gorm:"size:25; not null;unique;index" validate:"required"`
	Email    string  `json:"email" gorm:"size:100; not null;unique;index" validate:"required,email"`
	Profile  Profile `json:"profile" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ID"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(user User) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func GetUsers(c *fiber.Ctx) error {
	db := db.DBConn
	var users []User

	db.Preload("Profile").Find(&users)
	return c.JSON(users)
}

func CreateUser(c *fiber.Ctx) error {
	db := db.DBConn
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := ValidateStruct(*user)
	if errors != nil {
		c.JSON(errors)
		return nil
	}

	db.Preload("Profile").Create(&user)
	if user.ID > 0 && user.Profile.ID > 0 {
		return c.JSON(user)
	} else {
		return c.Status(400).Send([]byte("Bad Request!"))
	}
}

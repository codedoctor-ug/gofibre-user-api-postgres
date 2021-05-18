package models

import (
	"time"

	"github.com/abiiranathan/gofibre-gorm-relations/db"
	"github.com/gofiber/fiber/v2"
)

type Profile struct {
	ID        int       `json:"id" gorm:"not null;primaryKey"`
	FullName  string    `json:"name" gorm:"size:50;not null"`
	Mobile    string    `json:"mobile" gorm:"size:10;not null;unique"`
	BirthDate time.Time `json:"birth_date" gorm:"not null"`
}

type User struct {
	ID       int     `json:"id" gorm:"not null;primaryKey"`
	Username string  `json:"username" gorm:"size:25; not null;unique;index"`
	Email    string  `json:"email" gorm:"size:100; not null;unique;index"`
	Profile  Profile `json:"profile" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ID"`
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
		c.Status(400).Send([]byte(err.Error()))
	}

	db.Preload("Profile").Create(&user)
	if user.ID > 0 && user.Profile.ID > 0 {
		return c.JSON(user)
	} else {
		return c.Status(400).Send([]byte("Bad Request!"))
	}
}

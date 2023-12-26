package seeders

import (
	"log"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/models"
	"gorm.io/gorm"
)

type (
	UserSeeder interface {
		Seed() error
		Clear() error
	}
	userSeeder struct {
		db *gorm.DB
	}
)

func NewUserSeeder(db *gorm.DB) UserSeeder {
	return userSeeder{db: db}
}

// Implement seed method
func (s userSeeder) Seed() error {
	log.Println("UserSeeder running...")

	user := models.User{
		FirstName: "Super",
		LastName:  "Administrator",
		Email:     "superadmin@stream.co.th",
	}

	result := s.db.Create(&user)
	log.Println("UserSeeder seeded!")

	return result.Error
}

// Implement clear method
func (s userSeeder) Clear() error {
	log.Println("Clear UserSeeder...")
	result := s.db.Delete(&models.User{})
	log.Println("UserSeeder cleared!")

	return result.Error
}

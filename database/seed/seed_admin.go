package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/rizkisundara/project-management-api/config"
	"github.com/rizkisundara/project-management-api/models"
	"github.com/rizkisundara/project-management-api/utils"
)

func SeedAdmin() {
	password, _ := utils.HashPassword("admin123")

	admin := models.User{
		Name:     "Super admin",
		Email:    "admin@example.com",
		Password: password,
		Role:     "admin",
		PublicID: uuid.New(),
	}

	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error; err != nil {
		log.Println("Failed too seed admin", err)
	} else {
		log.Println("Admin seeded successfully")
	}
}

package repositories

import (
	"github.com/rizkisundara/project-management-api/config"
	"github.com/rizkisundara/project-management-api/models"
	"github.com/rizkisundara/project-management-api/utils"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
	FindAllUsersWithPagination(q utils.PaginationQuery) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := config.DB.Where("email = ?", email).First(&user).Error

	return &user, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id = ?", publicID).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAllUsersWithPagination(q utils.PaginationQuery) ([]models.User, int64, error) {
	var users []models.User

	baseQuery := utils.PaginationConfig{
		SearchFields: []string{
			"name",
			"email",
		},
		AllowedSortFields: map[string]string{
			"id":         "internal_id",
			"name":       "name",
			"email":      "email",
			"created_at": "created_at",
			"updated_at": "updated_at",
		},
		DefaultSortField: "internal_id",
		DefaultSortOrder: "DESC",
		UseILIKE:         true,
	}

	total, err := utils.CountWithFilter(config.DB.Model(&models.User{}), q, baseQuery)
	if err != nil {
		return nil, 0, err
	}

	query := utils.ApplyPagination(config.DB.Model(&models.User{}), q, baseQuery)
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Model(&models.User{}).
		Where("public_id = ?", user.PublicID).Updates(map[string]interface{}{
		"name": user.Name,
	}).Error
}

func (r *userRepository) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}

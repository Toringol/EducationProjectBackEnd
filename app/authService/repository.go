package authService

import "github.com/Toringol/EducationProjectBackEnd/app/models"

// UserRepository - interface for interaction with DB
type UserRepository interface {
	SelectUserByID(int64) (*models.User, error)
	SelectUserByEmail(string) (*models.User, error)
	CreateUser(*models.User) (int64, error)
}

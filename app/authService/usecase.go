package authService

import "github.com/Toringol/EducationProjectBackEnd/app/models"

// UserUsecase - interface for interaction User
type UserUsecase interface {
	SelectUserByID(int64) (*models.User, error)
	SelectUserByEmail(string) (*models.User, error)
	CreateUser(*models.User) (int64, error)
}

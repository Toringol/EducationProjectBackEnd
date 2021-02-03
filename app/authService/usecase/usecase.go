package usecase

import (
	"github.com/Toringol/EducationProjectBackEnd/app/authService"
	"github.com/Toringol/EducationProjectBackEnd/app/models"
)

// userUsecase - connector database and user
type userUsecase struct {
	repo authService.UserRepository
}

// NewUserUsecase - create new userUsecase structure
func NewUserUsecase(userRepo authService.UserRepository) authService.UserUsecase {
	return userUsecase{repo: userRepo}
}

// SelectUserByID - select user data by ID
func (us userUsecase) SelectUserByID(id int64) (*models.User, error) {
	return us.repo.SelectUserByID(id)
}

// SelectUserByEmail - select user data by Email
func (us userUsecase) SelectUserByEmail(email string) (*models.User, error) {
	return us.repo.SelectUserByEmail(email)
}

// CreateUser - create new user record in DB
func (us userUsecase) CreateUser(user *models.User) (int64, error) {
	return us.repo.CreateUser(user)
}

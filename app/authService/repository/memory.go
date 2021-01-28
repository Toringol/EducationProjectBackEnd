package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Toringol/EducationProjectBackEnd/app/authService"
	"github.com/Toringol/EducationProjectBackEnd/app/models"
	"github.com/spf13/viper"
)

// userRepository - structure that connects to DB and implement methods
type userRepository struct {
	DB *sql.DB
}

// NewUserMemoryRepository - create connection to DB
func NewUserMemoryRepository() authService.UserRepository {
	host := viper.GetString("DBHost")
	port := viper.GetInt("DBPort")
	user := viper.GetString("DBUser")
	password := viper.GetString("DBPassword")
	dbname := viper.GetString("DBName")

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Println(err)
		return nil
	}
	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil
	}

	return &userRepository{
		DB: db,
	}
}

// SelectUserByID - return user data by ID
func (repo *userRepository) SelectUserByID(id int64) (*models.User, error) {
	record := &models.User{}

	err := repo.DB.
		QueryRow("SELECT id, email, name, password, avatar, role WHERE id = ?", id).
		Scan(&record.ID, &record.Email, &record.Name, &record.Password, &record.Password, &record.Role)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// SelectUserByEmail - return user data by email
func (repo *userRepository) SelectUserByEmail(email string) (*models.User, error) {
	record := &models.User{}

	err := repo.DB.
		QueryRow("SELECT id, email, name, password, avatar, role WHERE email = ?", email).
		Scan(&record.ID, &record.Email, &record.Name, &record.Password, &record.Password, &record.Role)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// CreateUser - create new user record in DB
func (repo *userRepository) CreateUser(user *models.User) (int64, error) {
	result, err := repo.DB.Exec(
		"INSERT INTO users (`email`, `name`, `password`, `avatar`, `role`) VALUES (?, ?, ?, ?, ?)",
		user.Email,
		user.Name,
		user.Password,
		user.Avatar,
		user.Role,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

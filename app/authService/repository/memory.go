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
		QueryRow("SELECT id, email, name, password, avatar, role WHERE id = $1", id).
		Scan(&record.ID, &record.Email, &record.Name, &record.Password, &record.Avatar, &record.Role)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// SelectUserByEmail - return user data by email
func (repo *userRepository) SelectUserByEmail(email string) (*models.User, error) {
	record := &models.User{}

	err := repo.DB.
		QueryRow("SELECT id, email, name, password, avatar, role from users WHERE email = $1", email).
		Scan(&record.ID, &record.Email, &record.Name, &record.Password, &record.Avatar, &record.Role)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// CreateUser - create new user record in DB
func (repo *userRepository) CreateUser(user *models.User) (int64, error) {
	var id int64
	err := repo.DB.QueryRow(
		"INSERT INTO users (email, name, password, avatar, role) "+
			"VALUES ($1, $2, $3, $4, $5) "+
			"RETURNING id",
		user.Email,
		user.Name,
		user.Password,
		user.Avatar,
		user.Role,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

package main

import (
	"log"
	"net/http"

	authService "github.com/Toringol/EducationProjectBackEnd/app/authService/delivery/http"
	"github.com/Toringol/EducationProjectBackEnd/app/authService/repository"
	"github.com/Toringol/EducationProjectBackEnd/app/authService/usecase"
	"github.com/Toringol/EducationProjectBackEnd/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	authServiceListenAddr := viper.GetString("authServiceListenAddr")

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${method}] ${remote_ip}, ${uri} ${status} 'error':'${error}'\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{viper.GetString("frontEndAddr")},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	authService.NewAuthHandlers(e, usecase.NewUserUsecase(repository.NewUserMemoryRepository()))

	e.Logger.Fatal(e.Start(authServiceListenAddr))
}

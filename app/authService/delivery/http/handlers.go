package http

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"strconv"

	"github.com/Toringol/EducationProjectBackEnd/app/authService"
	"github.com/Toringol/EducationProjectBackEnd/app/models"
	"github.com/Toringol/EducationProjectBackEnd/app/sessionService/session"
	"github.com/Toringol/EducationProjectBackEnd/tools"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type authHandlers struct {
	usecase       authService.UserUsecase
	sessionClient session.SessionCheckerClient
}

// NewAuthHandlers - create user handlers
func NewAuthHandlers(e *echo.Echo, us authService.UserUsecase, sc session.SessionCheckerClient) {
	handlers := authHandlers{usecase: us, sessionClient: sc}

	e.POST("/login/", handlers.handleLogIn)
	e.POST("/signup/", handlers.handleSignUp)
}

func (uh *authHandlers) handleLogIn(ctx echo.Context) error {
	authCredentials := new(models.User)
	if err := ctx.Bind(authCredentials); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	userRecord, err := uh.usecase.SelectUserByEmail(authCredentials.Email)
	switch {
	case err == sql.ErrNoRows:
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect email or password")
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	oldPassDecrypted, err := base64.RawStdEncoding.DecodeString(userRecord.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	if !tools.CheckPass(oldPassDecrypted, authCredentials.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect email or password")
	}

	sessionID, err := uh.sessionClient.Create(ctx.Request().Context(), &session.Session{
		UserID:    strconv.FormatInt(userRecord.ID, 10),
		UserAgent: ctx.Request().UserAgent(),
		UserRole:  strconv.Itoa(userRecord.Role),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	log.Println(sessionID)

	userRecord.ID = 0
	userRecord.Password = ""

	return ctx.JSON(http.StatusOK, userRecord)
}

func (uh *authHandlers) handleSignUp(ctx echo.Context) error {
	userCredentials := new(models.User)
	if err := ctx.Bind(userCredentials); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bad request")
	}

	_, err := uh.usecase.SelectUserByEmail(userCredentials.Email)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	default:
		return echo.NewHTTPError(http.StatusConflict, "This email already registered")
	}

	userCredentials.Password = base64.RawStdEncoding.EncodeToString(
		tools.ConvertPass(userCredentials.Password))
	userCredentials.Avatar = viper.GetString("staticStoragePath") + "avatars/defaultAvatar.svg"

	_, err = uh.usecase.CreateUser(userCredentials)
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusInternalServerError, "Internal error")
	}

	userCredentials.ID = 0
	userCredentials.Password = ""

	return ctx.JSON(http.StatusOK, userCredentials)
}

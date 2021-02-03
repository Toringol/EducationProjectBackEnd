package tools

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
)

var cookieManager = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func SetCookie(ctx echo.Context, sessID string) error {
	value := map[string]string{
		"session_id": sessID,
	}

	encoded, err := cookieManager.Encode("session_token", value)
	if err != nil {
		return err
	}

	expiration := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   encoded,
		Path:    "/",
		Expires: expiration,
	}

	ctx.SetCookie(cookie)

	return nil
}

func ReadCookie(ctx echo.Context) (string, error) {
	cookie, err := ctx.Cookie("session_token")
	if err != nil {
		return "", err
	}

	value := make(map[string]string)
	err = cookieManager.Decode("session_token", cookie.Value, &value)
	if err != nil {
		return "", err
	}

	return value["session_id"], nil
}

func ClearCookie(ctx echo.Context) {
	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	ctx.SetCookie(cookie)
}

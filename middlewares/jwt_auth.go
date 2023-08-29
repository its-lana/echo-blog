package middlewares

import (
	"echo-blog/helper"
	"echo-blog/models"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type MyCustomClaims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}

func CreateToken(userId int) (string, error) {
	claims := MyCustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), //Token expires after 1 hour
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func validateToken(encodedToken string) (int, error) {
	signatureKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(encodedToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signatureKey, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if ok && token.Valid {
		return claims.UserId, nil
	} else {
		return 0, errors.New("token invalid")
	}

}

func UserAuthMiddlewares() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return helper.WrapResponse(http.StatusUnauthorized, "You are not Authorized!", &models.User{}).WriteToResponseBody(c.Response())
			}

			token := strings.Split(authHeader, " ")[1]
			userId, e := validateToken(token)
			if userId == 0 || e != nil {
				return helper.WrapResponse(http.StatusUnauthorized, "You are not Authorized!", &models.User{}).WriteToResponseBody(c.Response())
			}

			c.Set("userId", userId)
			return next(c)
		}
	}
}

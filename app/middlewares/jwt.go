package middlewares

import (
	"echo-fullstack/app/repository"
	"echo-fullstack/app/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get("Authorization")
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) != 2 {
				return c.JSON(http.StatusUnauthorized, utils.NewUnauthorizedError("Token Autorization salah"))
			}

			tokenStr := bearerToken[1]

			fmt.Println(tokenStr)

			UserID, err := repository.ValidateToken(tokenStr)
			if err != nil {
				fmt.Println("validasi token", err)
				return c.JSON(
					http.StatusUnauthorized,
					utils.NewUnauthorizedError(err.Error()),
				)
			}
			// fmt.Println("aut id", UserID)
			c.Set("user_id", UserID)
			c.Set("userId", UserID)

			return next(c)
		}
	}
}

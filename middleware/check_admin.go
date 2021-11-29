package middleware

import (
	"github.com/labstack/echo"
	"net/http"
)

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// handler logic
			if true {
				// tra ve loi
				return ctx.JSON(http.StatusInternalServerError, echo.Map{
					"message": "Bạn chưa đăng nhập",
				})
			}

			return next(ctx)
		}
	}
}

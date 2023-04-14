package handler

import (
	"fmt"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellarcd/internal/util/echttputil"
)

type ReflectorHandler func()

func HandleReflector() echo.HandlerFunc {
	return func(c echo.Context) error {
		resp, _ := io.ReadAll(c.Request().Body)
		fmt.Printf("%s %s\n", time.Now().Format(time.RFC822), string(resp))

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

package handler

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellarcd/internal/util/echttputil"
)

type ReflectorHandler func()

func HandleReflector() echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := make(map[string]interface{})
		err := json.NewDecoder(c.Request().Body).Decode(&resp)
		fmt.Printf("%#+v\n", resp)
		fmt.Printf("%#+v\n", err)

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

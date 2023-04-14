package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/go-backend/internal/service"
	"github.com/nmluci/go-backend/internal/util/echttputil"
	"github.com/nmluci/go-backend/pkg/errs"
)

func AuthorizationMiddleware(svc service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header

			var err error
			ctx := c.Request().Context()
			if token := header.Get("authorization"); token != "" {
				ctx, err = svc.AuthenticateSession(ctx, token)
			} else if token := header.Get("st-kagi"); token != "" {
				ctx, err = svc.AuthenticateService(ctx, token)
			} else {
				return echttputil.WriteErrorResponse(c, errs.ErrNoAccess)
			}

			if err != nil {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			c.SetRequest(c.Request().Clone(ctx))
			return next(c)
		}
	}
}

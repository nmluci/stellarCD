package handler

import (
	"context"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellarcd/internal/util/echttputil"
	"github.com/nmluci/stellarcd/pkg/dto"
	"github.com/rs/zerolog/log"
)

type DeploymentHandler func(context.Context, *dto.WebhoookRequest) error

func HandleDeployment(handler DeploymentHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &dto.WebhoookRequest{
			JobID: c.Param("jobID"),
		}

		if err := json.NewDecoder(c.Request().Body).Decode(&req.Webhook); err != nil {
			log.Printf("%#+v", err)
			return echttputil.WriteErrorResponse(c, err)
		}

		err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

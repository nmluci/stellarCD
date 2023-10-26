package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellarcd/internal/util/echttputil"
	"github.com/nmluci/stellarcd/pkg/dto"
)

type DeploymentHandler func(context.Context, *dto.WebhoookRequest) error

func HandleDeployment(handler DeploymentHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &dto.WebhoookRequest{
			JobID:     c.Param("jobID"),
			HeaderMap: c.Request().Header,
		}

		raw, _ := io.ReadAll(c.Request().Body)

		if err := json.Unmarshal(raw, &req.Webhook); err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		dst := &bytes.Buffer{}
		json.Compact(dst, raw)
		req.RawBody = dst.Bytes()

		err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

type SimpleDeploymentHandler func(context.Context, *dto.WebhoookRequest) error

func HandleSimpleDeployment(handler DeploymentHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &dto.WebhoookRequest{
			JobID:     c.Param("jobID"),
			HeaderMap: c.Request().Header,
			Webhook: map[string]interface{}{
				"version": c.QueryParam("version"),
			},
		}

		err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

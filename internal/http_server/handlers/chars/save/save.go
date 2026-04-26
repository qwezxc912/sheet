package charsave

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	resp "github.com/qweq1232/dnd_form/internal/domane/api/response"
)

type Request struct {
	UserID int32           `json:"user_id"`
	Stats  json.RawMessage `json:"stats"`
}

type Response struct {
	resp.Response
}

type CharSaver interface {
	SaveChar(ctx context.Context, stats []byte, userID int32) error
}

func New(ctx context.Context, log *slog.Logger, cs CharSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.chars.save.New"

		log = log.With(slog.String("op", op))

		var req Request

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))

			c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := cs.SaveChar(ctx, req.Stats, req.UserID); err != nil {
			log.Error("failed to save character", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		log.Info("character successfuly saved")

		c.JSON(http.StatusOK, resp.OK())
	}
}

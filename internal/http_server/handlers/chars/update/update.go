package charupdate

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	resp "github.com/qweq1232/dnd_form/internal/domane/api/response"
)

type Request struct {
	Stats json.RawMessage `json:"stats"`
}

type Response struct {
	resp.Response
}

type CharUpdater interface {
	UpdateChar(ctx context.Context, stats []byte, id int32) error
}

func New(ctx context.Context, log *slog.Logger, cu CharUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.chars.update.New"

		log = log.With(slog.String("op", op))

		var req Request

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))

			c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))

			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error("failed to parse id from path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid id"))

			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := cu.UpdateChar(ctx, req.Stats, int32(id)); err != nil {
			log.Error("failed to update character", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		log.Info("character successfuly updated")

		c.JSON(http.StatusOK, resp.OK())
	}
}

package chardelete

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	resp "github.com/qweq1232/dnd_form/internal/domane/api/response"
)

type Response struct {
	resp.Response
}

type CharDeleter interface {
	DeleteChar(ctx context.Context, id int32) error
}

func New(ctx context.Context, log *slog.Logger, cd CharDeleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.chars.delete.New"

		log = log.With(slog.String("op", op))

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error("failed to parse id from path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid id"))

			return
		}

		if err := cd.DeleteChar(ctx, int32(id)); err != nil {
			log.Error("failed to delete character", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		log.Info("character successfuly deleted")

		c.JSON(http.StatusOK, resp.OK())
	}
}

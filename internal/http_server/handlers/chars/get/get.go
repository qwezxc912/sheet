package get

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	resp "github.com/qweq1232/dnd_form/internal/domane/api/response"
	models "github.com/qweq1232/dnd_form/internal/domane/models/char"
)

type Response struct {
	resp.Response
	char models.Char `json:"stats"`
}

type CharGetter interface {
	GetChar(ctx context.Context, id, userID int32) (*models.Char, error)
}

func New(ctx context.Context, log *slog.Logger, cg CharGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.chars.get.New"

		log = log.With(slog.String("op", op))

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Error("failed to parse id from path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid id"))

			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			log.Error("failed to parse user_id from path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid id"))

			return
		}

		char, err := cg.GetChar(ctx, int32(id), int32(userID))
		if err != nil {
			log.Error("failed to get stats", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		log.Info("get character successfuly", slog.Any("char id", id))

		c.JSON(http.StatusOK, Response{resp.OK(), *char})
	}
}

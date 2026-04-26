package getall

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
	Chars []models.Char `json:"chars"`
}

type AllCharGetter interface {
	GetAllChar(ctx context.Context, userID int32) ([]models.Char, error)
}

func New(ctx context.Context, log *slog.Logger, gac AllCharGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.chars.get_all.New"

		log = log.With(slog.String("op", op))

		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			log.Error("failed to parse path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid user_id"))

			return
		}

		chars, err := gac.GetAllChar(ctx, int32(id))
		if err != nil {
			log.Error("failed to get all characters", slog.Any("err", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		log.Info("all characters get seccessfuly", slog.Any("user id", id))

		c.JSON(http.StatusOK, Response{resp.OK(), chars})
	}
}

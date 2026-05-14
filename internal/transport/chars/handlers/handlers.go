package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	resp "github.com/qweq1232/dnd_form/internal/domane/api/response"
	models "github.com/qweq1232/dnd_form/internal/domane/models/char"
	"github.com/qweq1232/dnd_form/internal/service"
)

type handler struct {
	service *service.Service
	log     *slog.Logger
}

type CreateRequest struct {
	UserID int32           `json:"user_id"`
	Stats  json.RawMessage `json:"stats"`
}

type UpdateRequest struct {
	Stats json.RawMessage `json:"stats"`
}

func New(serv *service.Service, log *slog.Logger) *handler {
	return &handler{service: serv, log: log}
}

func (h handler) Create(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "chars.handlers.Create"

		h.log = h.log.With(slog.String("op", op))

		var req CreateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			h.log.Error("failed to decode request body", slog.Any("error", err))

			c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))

			return
		}

		h.log.Info("request body decoded", slog.Any("req", req))

		char := models.Char{UserID: req.UserID, Stats: req.Stats}
		id, err := h.service.CreateChar(ctx, char)
		if err != nil {
			h.log.Error("failed to save character", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		h.log.Info("character successfuly saved")

		c.JSON(http.StatusOK, gin.H{"status": resp.OK(), "id": id})
	}
}

func (h handler) Get(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "chars.handlers.Get"

		h.log = h.log.With(slog.String("op", op))

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.log.Error("failed to parse id from path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid id"))

			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			h.log.Error("failed to parse user_id from path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid id"))

			return
		}

		char, err := h.service.Char(ctx, int32(id), int32(userID))
		if err != nil {
			h.log.Error("failed to get stats", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		h.log.Info("get character successfuly", slog.Any("char id", id))

		c.JSON(http.StatusOK, char)
	}
}

func (h handler) GetAll(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "chars.handlers.GetAll"

		h.log = h.log.With(slog.String("op", op))

		id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			h.log.Error("failed to parse path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid user_id"))

			return
		}

		chars, err := h.service.AllChar(ctx, int32(id))
		if err != nil {
			h.log.Error("failed to get all characters", slog.Any("err", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		h.log.Info("all characters get seccessfuly", slog.Any("user id", id))

		c.JSON(http.StatusOK, chars)
	}
}

func (h handler) Update(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "chars.handlers.Update"

		h.log = h.log.With(slog.String("op", op))

		var req UpdateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			h.log.Error("failed to decode request body", slog.Any("error", err))

			c.JSON(http.StatusBadRequest, resp.Error("failed to decode request"))

			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.log.Error("failed to parse id from path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid id"))

			return
		}

		h.log.Info("request body decoded", slog.Any("req", req))

		char := models.Char{ID: int32(id), Stats: req.Stats}

		if err := h.service.UpdateChar(ctx, char); err != nil {
			h.log.Error("failed to update character", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		h.log.Info("character successfuly updated")

		c.JSON(http.StatusOK, resp.OK())
	}
}

func (h handler) Delete(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "chars.handlers.Delete"

		h.log = h.log.With(slog.String("op", op))

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.log.Error("failed to parse id from path param", slog.Any("err", err))

			c.JSON(http.StatusBadRequest, resp.Error("invalid id"))

			return
		}

		if err := h.service.DeleteChar(ctx, int32(id)); err != nil {
			h.log.Error("failed to delete character", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, resp.Error("internal error"))

			return
		}

		h.log.Info("character successfuly deleted")

		c.JSON(http.StatusOK, resp.OK())
	}
}

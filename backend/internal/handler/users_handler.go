package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	dbq "github.com/valentineejk/piple/db/sqlc"
	"github.com/valentineejk/piple/internal/helpers"
	"github.com/valentineejk/piple/internal/model"
)

const pgUniqueViolation = "23505"

func (h *Handler) CreateUser(c *gin.Context) {

	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid body data",
		})
		return
	}

	if !model.ValidUserRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid role",
			"message": "role must be one of: employee, procurement, ceo, admin",
		})
		return
	}

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to hash password",
		})
		return
	}

	user, err := h.queries.CreateUser(c.Request.Context(), dbq.CreateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashed,
		Role:      req.Role,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			c.JSON(http.StatusConflict, gin.H{
				"error":   pgErr.Message,
				"message": "a user with this email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    userResponseFromUser(user),
		"message": "user created successfully",
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := helpers.ParseUUID(c.Param("id"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "invalid body data",
		})
		return
	}

	if req.Role != nil && !model.ValidUserRoles[*req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid role",
			"message": "role must be one of: employee, procurement, ceo, admin",
		})
		return
	}

	user, err := h.queries.UpdateUser(c.Request.Context(), dbq.UpdateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
		ID:        id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			c.JSON(http.StatusConflict, gin.H{
				"error":   pgErr.Message,
				"message": "a user with this email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    userResponseFromUser(user),
		"message": "user updated successfully",
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := helpers.ParseUUID(c.Param("id"))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := h.queries.SoftDeleteUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found or already deleted",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    userResponseFromUser(user),
		"message": "user deleted successfully",
	})
}

func userResponseFromUser(u dbq.User) model.UserWriteResponse {
	return model.UserWriteResponse{
		ID:        helpers.UUIDToString(u.ID),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Time,
		DeletedAt: helpers.TimestampToTimePtr(u.DeletedAt),
	}
}

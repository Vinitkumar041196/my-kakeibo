package http

import (
	"kakeibo-api/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary Create new user
// @Description Creates new system user
// @Tags User
// @Produce json
// @Param req body domain.SignUpRequest true "user details"
// @Success 200 {object} domain.UserSuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /signup [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	req := domain.SignUpRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	err = h.userService.CreateUser(&req.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.UserSuccessResponse{Message: "success"})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user details
// @Tags User
// @Produce json
// @Param user_id path string true "user details"
// @Param req body domain.SignUpRequest true "user details"
// @Success 200 {object} domain.UserSuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /user/{user_id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "invalid user id"})
		return
	}

	if userId == 0 {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "invalid user id"})
		return
	}

	req := domain.User{}
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	err = h.userService.UpdateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.UserSuccessResponse{Message: "success"})
}

// GetUser godoc
// @Summary Get user
// @Description Get user details
// @Tags User
// @Produce json
// @Param user_id path string true "user details TODO:to be changed to get the user id from token"
// @Success 200 {object} domain.User
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /user/{user_id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userId := c.Param("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "invalid user id"})
		return
	}

	user, err := h.userService.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user
// @Tags User
// @Produce json
// @Param user_id path string true "user details TODO:to be changed to get the user id from token"
// @Success 200 {object} domain.UserSuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /user/{user_id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userId := c.Param("user_id")
	if userId != "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "invalid user id"})
		return
	}

	err := h.userService.DeleteUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.UserSuccessResponse{Message: "success"})
}

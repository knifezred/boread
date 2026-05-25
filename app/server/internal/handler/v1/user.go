package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"boread/internal/dto"
	"boread/internal/model"
	"boread/internal/service"
	jwtPkg "boread/pkg/jwt"
	"boread/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 3001, err.Error())
		return
	}

	response.Success(c, toUserResponse(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}

	user, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 2001, err.Error())
		return
	}

	token, err := jwtPkg.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Error(c, 5001, "failed to generate token")
		return
	}

	response.Success(c, dto.LoginResponse{
		Token:    token,
		UserInfo: toUserResponse(user),
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 3002, "user not found")
		return
	}

	response.Success(c, toUserResponse(user))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}

	user, err := h.userService.Update(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, 3003, err.Error())
		return
	}

	response.Success(c, toUserResponse(user))
}

func (h *UserHandler) List(c *gin.Context) {
	var req dto.PageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 1001, err.Error())
		return
	}

	users, total, err := h.userService.List(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 5002, err.Error())
		return
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, toUserResponse(&u))
	}

	response.Success(c, dto.PageResponse{
		List:     userResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, 1001, "invalid id")
		return
	}

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, 3004, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, 3002, "user not found")
		return
	}

	response.Success(c, toUserResponse(user))
}

func toUserResponse(u *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Email:    u.Email,
		Phone:    u.Phone,
		Avatar:   u.Avatar,
		Status:   u.Status,
	}
}

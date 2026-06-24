package handler

import (
	"studyspot/internal/domain"
	"studyspot/internal/service"
	"studyspot/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	user, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			response.BadRequest(c, "Пользователь с таким email уже существует")
			return
		}
		response.InternalError(c, "Ошибка регистрации: "+err.Error())
		return
	}

	response.Created(c, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Неверный формат запроса: "+err.Error())
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			response.Unauthorized(c, "Неверный email или пароль")
			return
		}
		response.InternalError(c, "Ошибка авторизации: "+err.Error())
		return
	}

	response.Success(c, domain.LoginResponse{
		Token: token,
		User:  *user,
	})
}

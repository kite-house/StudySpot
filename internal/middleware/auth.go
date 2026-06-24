package middleware

import (
	"strings"

	"studyspot/pkg/jwt"
	"studyspot/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Требуется авторизация")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Неверный формат токена")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwt.ValidateToken(token, jwtSecret)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				response.Unauthorized(c, "Срок действия токена истёк")
			} else {
				response.Unauthorized(c, "Неверный токен")
			}
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID.String())
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "Доступ запрещён")
			c.Abort()
			return
		}

		if role != "admin" {
			response.Forbidden(c, "Требуются права администратора")
			c.Abort()
			return
		}

		c.Next()
	}
}

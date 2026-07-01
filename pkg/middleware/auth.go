package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kushaljangra/e-commerce/pkg/auth"
	"github.com/kushaljangra/e-commerce/pkg/contextkeys"
)

func tokenFromRequest(c *gin.Context) string {
	if authHeader := c.GetHeader("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	authCookie, err := c.Cookie("token")
	if err != nil || authCookie == "" {
		return ""
	}

	return authCookie
}

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		encodedToken := tokenFromRequest(c)
		if encodedToken == "" {
			c.Set("userID", "")
			c.Next()
			return
		}

		token, err := auth.ValidateToken(encodedToken)
		if err != nil {
			c.Set("userID", "")
			c.Next()
			return
		}

		if claims, ok := token.Claims.(*auth.JWTCustomClaims); ok && token.Valid {
			c.Set("userID", claims.UserID)
			ctxWithVal := context.WithValue(c.Request.Context(), contextkeys.UserIDKey, claims.UserID)

			c.Request = c.Request.WithContext(ctxWithVal)
		} else {
			c.Set("userID", "")
		}

		c.Next()
	}
}

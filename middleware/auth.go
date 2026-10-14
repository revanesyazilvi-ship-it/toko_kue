package middleware

import (
	"net/http"
	"strings"

	"toko_kue/policy"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte("secret_key_toko_kue_nesya_123")

func AuthMiddleware(secretKey string) gin.HandlerFunc {
		return func(c *gin.Context) {
				authHeader := c.GetHeader("Authorization")
				if authHeader == "" {
						c.JSON(http.StatusUnauthorized, gin.H{"error": "Header otorrisasi dibutuhkan"})
						c.Abort()
						return
				}

				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				tokenString = strings.TrimSpace(tokenString)
				token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
						return []byte(secretKey),nil
				})
			if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kadaluarsa"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Klaim token tidak valid"})
			c.Abort()
			return
		}

		c.Set(policy.UserCtxKey, claims["user_id"])
		c.Set(policy.RoleCtxKey, claims["role"])

		c.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(policy.RoleCtxKey)
		if !exists || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: Kamu tidak memiliki izin"})
			c.Abort()
			return
		}
		c.Next()
	}
	
		}
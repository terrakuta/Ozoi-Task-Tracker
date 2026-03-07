package middleware

import (
	"Ozoi/internal/config"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("invalid token algo: %s", token.Method.Alg())
			}

			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})

			if token == nil { // logging unusual cases
				log.Printf("Malformed token from %s: %v", c.ClientIP(), err)
			} else if err != nil {
				log.Printf("Invalid token from %s: %v (alg: %v)",
					c.ClientIP(), err, token.Header["alg"])
			} else if !token.Valid {
				log.Printf("Token not valid from %s (expired or bad signature)", c.ClientIP())
			}

			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		expRaw, exists := claims["exp"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing expiration (exp claim required)"})
			c.Abort()
			return
		}

		exp, ok := expRaw.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid expiration claim format"})
			c.Abort()
			return
		}

		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

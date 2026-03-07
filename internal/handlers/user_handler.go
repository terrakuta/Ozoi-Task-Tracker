package handlers

import (
	"Ozoi/internal/config"
	"Ozoi/internal/models"
	"Ozoi/repository"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRespone struct {
	Token string `json:"token"`
}

// CreateUserHandler godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      RegisterRequest  true  "Email and password"
// @Success      201    {object}  models.OzoiUser
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /auth/register [post]
func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error: ": err.Error()})
			return
		}

		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 character long"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error: ": "Failed to hash password " + err.Error()})
			return
		}

		user := &models.OzoiUser{
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		createdUser, err := repository.CreateUser(pool, user)

		if err != nil {
			if strings.Contains(err.Error(), "23505") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already registered"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, createdUser)
	}
}

// LoginHandler godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      LoginRequest      true  "Email and password"
// @Success      200    {object}  map[string]string  "message: Login successful"
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /auth/login [post]
func LoginHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByEmail(pool, req.Email)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(time.Hour).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token" + err.Error()})
			return
		}

		c.SetCookie("access_token", tokenString, 3600, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	}
}

// LogoutHandler godoc
// @Summary      Logout
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]string  "message: Logged out"
// @Router       /auth/logout [post]
func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("access_token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	}
}

// MeHandler godoc
// @Summary      Check current session
// @Tags         auth
// @Produce      json
// @Success      200
// @Failure      401  {object}  map[string]string
// @Security     CookieAuth
// @Router       /auth/me [get]
func MeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Status(http.StatusOK)
	}
}

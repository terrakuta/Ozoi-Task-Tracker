package main

import (
	_ "Ozoi/docs"
	"Ozoi/internal/config"
	"Ozoi/internal/database"
	"Ozoi/internal/handlers"
	"Ozoi/internal/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Package main Ozoi API
//
// @title           Ozoi API
// @version         1.0
// @description     Task manager REST API
//
// @host            localhost:8080
// @BasePath        /
//
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token
func main() {

	var cfg *config.Config
	var err error

	cfg, err = config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	swgCfg, err := config.LoadSwaggerConfig()
	if err != nil {
		log.Fatalf("failed to load swagger config: %v", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Println("Failed to connect DB")
	}
	defer pool.Close()

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)

	if err != nil {
		log.Fatalf("Cannot initialize migrations: ", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrations failed: ", err)
	}

	log.Println("Migrations have status success")

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	if swgCfg.SwaggerUser != "" && swgCfg.SwaggerPassword != "" {
		swaggerGroup := router.Group("/swagger", gin.BasicAuth(gin.Accounts{
			swgCfg.SwaggerUser: swgCfg.SwaggerPassword,
		}))

		swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.GET("/hello", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message":  "Hello World",
			"status":   "Success",
			"database": "Connected",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))
	router.POST("/auth/logout", handlers.LogoutHandler())

	protected := router.Group("/ozoi")
	auth := router.Group("/auth")

	protected.Use(middleware.AuthMiddleware(cfg))

	protected.POST("", handlers.CreateTaskHandler(pool))
	protected.GET("", handlers.GetAllTasksHandler(pool))
	protected.GET("/:id", handlers.GetTaskByIDHandler(pool))
	protected.PUT("/:id", handlers.UpdateTaskByIDHandler(pool))
	protected.DELETE("/:id", handlers.DeleteTaskByIDHandler(pool))

	auth.Use(middleware.AuthMiddleware(cfg))
	auth.GET("/me", handlers.MeHandler())

	router.Run(":" + cfg.Port)

}

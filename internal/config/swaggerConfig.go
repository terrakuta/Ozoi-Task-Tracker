package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SwaggerConfig struct {
	DatabaseURL     string
	Port            string
	JWTSecret       string
	SwaggerUser     string
	SwaggerPassword string
}

func LoadSwaggerConfig() (*SwaggerConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("env file not found, using system env")
	}

	return &SwaggerConfig{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		Port:            os.Getenv("PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		SwaggerUser:     os.Getenv("SWAGGER_USER"),
		SwaggerPassword: os.Getenv("SWAGGER_PASSWORD"),
	}, nil
}

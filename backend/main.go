package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	authhandler "webview-login/backend/auth/handler"
	authmodel "webview-login/backend/auth/model"
	"webview-login/backend/auth/repository"
	"webview-login/backend/auth/service"
	"webview-login/backend/config"
	"webview-login/backend/database"
	"webview-login/backend/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using environment variables")
	}

	cfg := config.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db := database.Connect(dsn, &authmodel.User{})

	// auth domain wiring
	userRepo := repository.NewUserRepo(db)
	registerHandler := authhandler.NewRegisterHandler(service.NewRegisterService(userRepo))
	loginHandler := authhandler.NewLoginHandler(service.NewLoginService(userRepo, cfg.JWTSecret))

	r := gin.Default()
	r.Use(middleware.CORS(cfg.FrontendOrigin))

	api := r.Group("/api")
	{
		api.POST("/register", registerHandler.Handle)
		api.POST("/login", loginHandler.Handle)
	}

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

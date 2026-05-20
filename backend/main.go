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
	"webview-login/backend/memory"
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
	db := database.Connect(dsn, &authmodel.UserInfo{})
	rdb := memory.Connect(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)

	// auth domain wiring
	userRepo := repository.NewDatabaseRepo(db)
	tokenRepo := repository.NewRedisRepo(rdb, cfg.RefreshTokenTTL, cfg.RedisRefreshTokenKey)

	registerHandler := authhandler.NewRegisterHandler(service.NewRegisterService(userRepo, cfg.EmailPepper, cfg.AESKey))
	loginHandler := authhandler.NewLoginHandler(service.NewLoginService(userRepo, tokenRepo, cfg.JWTSecret, cfg.EmailPepper))
	refreshHandler := authhandler.NewRefreshHandler(service.NewRefreshService(userRepo, tokenRepo, cfg.JWTSecret))
	logoutHandler := authhandler.NewLogoutHandler(service.NewLogoutService(tokenRepo))

	r := gin.Default()
	r.Use(middleware.CORS(cfg.FrontendOrigin))

	api := r.Group("/api")
	{
		api.POST("/register", registerHandler.Handle)
		api.POST("/login", loginHandler.Handle)
		api.POST("/auth/refresh", refreshHandler.Handle)
		api.POST("/auth/logout", logoutHandler.Handle)
	}

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

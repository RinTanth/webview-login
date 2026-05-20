package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	authhandler "webview-login/backend/auth/handler"
	authmodel "webview-login/backend/auth/model"
	"webview-login/backend/auth/repository"
	"webview-login/backend/auth/service"
	"webview-login/backend/database"
	"webview-login/backend/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using environment variables")
	}

	db := database.Connect(getenv("DB_PATH", "app.db"), &authmodel.User{})

	// auth domain wiring
	userRepo := repository.NewUserRepo(db)
	registerHandler := authhandler.NewRegisterHandler(service.NewRegisterService(userRepo))
	loginHandler := authhandler.NewLoginHandler(service.NewLoginService(userRepo))

	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.POST("/register", registerHandler.Handle)
		api.POST("/login", loginHandler.Handle)
	}

	port := getenv("PORT", "8080")
	log.Printf("server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

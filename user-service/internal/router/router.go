package router

import (
	"EcomMaster/pkg/logger"
	_ "user-service/docs"

	"user-service/internal/handler"

	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	userRepository := &repository.UserRepository{DB: db}
	userService := &service.UserService{Repo: *userRepository}
	userHandler := &handler.UserHandler{UserService: userService}

	r := gin.Default()
	r.Use(cors.Default())
	r.Use(middleware.CORSMiddleware())
	r.Use(logger.Logging())
	r.Use(middleware.ErrorHandler())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/auth")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/profile", userHandler.Profile)
		auth.PUT("/profile", userHandler.UpdateProfile)
	}

	return r
}

package server

import (
	"fmt"
	_ "kakeibo-api/docs"
	loginHttpDelivery "kakeibo-api/login/delivery/http"
	userHttpDelivery "kakeibo-api/user/delivery/http"
	"kakeibo-api/server/middleware"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (srv *Server) NewRouter() *gin.Engine {
	router := gin.Default()
	// setup cors
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"*"}
	cors_config.AllowMethods = []string{"*"}
	cors_config.AllowHeaders = []string{"*"}
	router.Use(cors.New(cors_config))
	
	router.Static("/web", "./static")

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK", "timestamp": time.Now(), "uptime": fmt.Sprintf("%v", time.Since(serverStart))})
	})

	loginHandler := loginHttpDelivery.NewLoginHandler(srv.LoginService)
	userHandler := userHttpDelivery.NewUserHandler(srv.UserService)
	
	v1 := router.Group("/v1")
	v1.POST("/signup", userHandler.CreateUser)
	v1.POST("/login", loginHandler.Login)
	v1.POST("/refreshtoken", loginHandler.RefreshToken)
	v1.POST("/forgotpassword/verifyemail", loginHandler.VerifyEmailAndSendOTP)
	v1.POST("/forgotpassword/verifyotp", loginHandler.ForgotPasswordVerifyOTP)
	v1.POST("/forgotpassword/reset", loginHandler.ForgotPassword)

	if os.Getenv("ENABLE_AUTH_MIDDLEWARE") == "true" {
		v1.Use(middleware.AuthenticationMiddleware())
	}

	return router
}


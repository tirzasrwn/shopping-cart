package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tirzasrwn/shopping-cart/internal/controllers/middleware"
	v1 "github.com/tirzasrwn/shopping-cart/internal/controllers/v1"
)

func InitializeRouter() (router *gin.Engine) {
	router = gin.Default()
	router.Use(
		middleware.CORSMiddleware(),
	)

	v1route := router.Group("/v1")
	{
		v1route.GET("/user", v1.GetUserByEmail)
	}

	return router
}

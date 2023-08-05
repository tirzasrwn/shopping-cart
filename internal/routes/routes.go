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
		v1route.POST("/user", v1.PostUser)
		v1route.GET("/category", v1.GetCategories)
		v1route.GET("/product", v1.GetProductByCategoryID)
		v1route.POST("/order", v1.InsertOrder)
	}

	return router
}

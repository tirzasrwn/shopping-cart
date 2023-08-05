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
		v1route.POST("/login", v1.Authenticate)
		v1route.GET("/category", v1.GetCategories)
		v1route.GET("/product", v1.GetProductByCategoryID)
		v1route.POST("/user", v1.PostUser)

		user := v1route.Group("/user")
		user.Use(middleware.AdminAuth.AuthRequest())
		{
			user.GET("", v1.GetUserInformation)
			user.GET("/order", v1.GetUserOrder)
			user.POST("/order", v1.InsertOrder)
			user.DELETE("/order/:order_id", v1.DeleteOrder)
			user.GET("/payment", v1.GetUserPayment)
			user.POST("/checkout", v1.Checkout)
		}
	}

	return router
}

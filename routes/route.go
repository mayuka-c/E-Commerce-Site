package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mayuka-c/e-commerce-site/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", controllers.SearchProducts())
	incomingRoutes.GET("/users/search", controllers.SearchProductsByQuery())
}

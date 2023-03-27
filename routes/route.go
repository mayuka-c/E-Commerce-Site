package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mayuka-c/e-commerce-site/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine, handler *controllers.Application) {
	incomingRoutes.POST("/users/signup", handler.SignUp())
	incomingRoutes.POST("/users/login", handler.Login())
	incomingRoutes.POST("/admin/addproduct", handler.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", handler.SearchProducts())
	incomingRoutes.GET("/users/search", handler.SearchProductsByQuery())
}

func ProductRoutes(incomingRoutes *gin.Engine, handler *controllers.Application) {
	incomingRoutes.POST("/addtocart", handler.AddToCart())
	incomingRoutes.DELETE("/removeitem", handler.RemoveItemFromCart())
	incomingRoutes.GET("/listcart", handler.GetItemFromCart())
	incomingRoutes.POST("/cartcheckout", handler.BuyFromCart())
	incomingRoutes.POST("/instantbuy", handler.InstantBuy())
}

func AddressRoutes(incomingRoutes *gin.Engine, handler *controllers.Application) {
	incomingRoutes.POST("/addaddress", handler.AddAddress())
	incomingRoutes.PUT("/edithomeaddress", handler.EditHomeAddress())
	incomingRoutes.PUT("/editworkaddress", handler.EditWorkAddress())
	incomingRoutes.DELETE("/deleteaddresses", handler.DeleteAddress())
}

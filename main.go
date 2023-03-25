package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mayuka-c/e-commerce-site/config"
	"github.com/mayuka-c/e-commerce-site/controllers"
	"github.com/mayuka-c/e-commerce-site/database"
	"github.com/mayuka-c/e-commerce-site/middleware"
	"github.com/mayuka-c/e-commerce-site/routes"
)

func main() {
	var serviceConfig config.ServiceConfig

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItemFromCart())
	router.GET("/listcart", app.GetItemFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.DELETE("/deleteaddresses", controllers.DeleteAddress())

	log.Println("E-commerce is running at port: ", serviceConfig.APIPort)
	log.Fatal(router.Run(":" + strconv.Itoa(serviceConfig.APIPort)))
}

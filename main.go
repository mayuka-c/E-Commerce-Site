package main

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/mayuka-c/e-commerce-site/config"
	"github.com/mayuka-c/e-commerce-site/controllers"
	"github.com/mayuka-c/e-commerce-site/database"
	"github.com/mayuka-c/e-commerce-site/middleware"
	"github.com/mayuka-c/e-commerce-site/routes"
	"github.com/mayuka-c/e-commerce-site/tokens"
)

var (
	ctx = context.Background()
)

var serviceConfig config.ServiceConfig
var dbConfig config.DBConfig

func init() {
	serviceConfig = config.GetServiceConfig(ctx)
	dbConfig = config.GetDBConfig(ctx)
}

func main() {

	dbClient := database.DBSet(dbConfig)
	tokenGenerator := tokens.NewTokenGenerator(dbClient)
	app := controllers.NewApplication(dbClient, tokenGenerator)

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router, app)
	router.Use(middleware.Authentication(tokenGenerator))

	routes.ProductRoutes(router, app)
	routes.AddressRoutes(router, app)

	log.Println("E-commerce is running at port: ", serviceConfig.APIPort)
	log.Fatal(router.Run(":" + strconv.Itoa(serviceConfig.APIPort)))
}

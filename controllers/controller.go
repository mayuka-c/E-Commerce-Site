package controllers

import (
	"github.com/mayuka-c/e-commerce-site/database"
	"github.com/mayuka-c/e-commerce-site/tokens"
)

type Application struct {
	dbClient    *database.DBClient
	tokenClient *tokens.TokenGenrator
}

func NewApplication(dbClient *database.DBClient, tokenClient *tokens.TokenGenrator) *Application {
	return &Application{
		dbClient:    dbClient,
		tokenClient: tokenClient,
	}
}

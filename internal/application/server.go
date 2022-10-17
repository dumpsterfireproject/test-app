package application

import (
	"context"
	"log"

	testapp "github.com/cycle-labs/test-app"
	"github.com/cycle-labs/test-app/internal/adapters/api"
	"github.com/cycle-labs/test-app/internal/adapters/inventorydb"
	"github.com/cycle-labs/test-app/internal/domain/ports"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	apiRoot = "/api"
)

func StartServer() {
	ctx := context.Background()
	db, err := CreateInMemoryDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repository := inventorydb.NewGengiInventoryStockRepository(db)
	err = repository.SeedData(ctx)
	if err != nil {
		log.Fatal(err)
	}
	stockService, _ := ports.NewInventoryStock(repository)

	authenticatedAPI := api.NewGinAuthenticatorEndpoint()
	jsonAPI := api.NewGinJsonAPIEndpoints(stockService)
	graphQLAPI := api.NewGraphQLEndpoints(stockService, apiRoot)

	ginEngine := gin.Default()
	authenticatedRoot := authenticatedAPI.AuthenticatedRouterGroup(ginEngine, apiRoot)
	jsonAPI.AddInventoryRoutes(authenticatedRoot)
	graphQLAPI.AddGraphQLRoutes(authenticatedRoot)

	fs := api.EmbedFolder(testapp.EmbeddedFiles, "webui/build", true)
	ginEngine.Use(static.Serve("/", fs))

	ginEngine.Run(":8080")
}

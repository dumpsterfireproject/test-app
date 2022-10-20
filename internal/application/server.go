package application

import (
	"context"
	"log"

	testapp "github.com/dumpsterfireproject/test-app"
	"github.com/dumpsterfireproject/test-app/internal/adapters/api"
	"github.com/dumpsterfireproject/test-app/internal/adapters/inventorydb"
	"github.com/dumpsterfireproject/test-app/internal/domain/ports"
	"github.com/gin-contrib/cors"
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
	ginEngine.Use(corsHandler())
	authenticatedRoot := authenticatedAPI.AuthenticatedRouterGroup(ginEngine, apiRoot)
	jsonAPI.AddInventoryRoutes(authenticatedRoot)
	graphQLAPI.AddGraphQLRoutes(authenticatedRoot)

	fs := api.EmbedFolder(testapp.EmbeddedFiles, "webui/build", true)
	ginEngine.Use(static.Serve("/", fs))

	ginEngine.Run(":8080")
}

// TODO: Update this to be less permissive
func corsHandler() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	return cors.New(config)
}

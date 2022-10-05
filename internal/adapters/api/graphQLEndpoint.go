package api

import (
	"fmt"

	"github.com/cycle-labs/test-app/internal/domain/ports"
	"github.com/gin-gonic/gin"
)

const (
	graphQLPath = "graphql"
)

type GraphQLEndpoints struct {
	service GraphQLService
}

func NewGraphQLEndpoints(stockService ports.InventoryStockService, root string) *GraphQLEndpoints {
	graphQLEndpoint := fmt.Sprintf("%s/%s", root, graphQLPath)
	api := &GraphQLEndpoints{
		service: NewGraphQLService(stockService, graphQLEndpoint),
	}
	return api
}

func (api *GraphQLEndpoints) AddGraphQLRoutes(root *gin.RouterGroup) {
	root.POST(graphQLPath, api.service.graphqlHandler())
	root.GET("graphiql", api.service.playgroundHandler())
}

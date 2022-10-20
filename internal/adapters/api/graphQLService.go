package api

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dumpsterfireproject/test-app/internal/adapters/graph"
	"github.com/dumpsterfireproject/test-app/internal/adapters/graph/generated"
	"github.com/dumpsterfireproject/test-app/internal/domain/ports"
	"github.com/gin-gonic/gin"
)

type GraphQLService interface {
	graphqlHandler() gin.HandlerFunc
	playgroundHandler() gin.HandlerFunc
}

func NewGraphQLService(stockService ports.InventoryStockService, graphQLEndpoint string) GraphQLService {
	return &GraphQLServiceImpl{stockService: stockService, graphQLEndpoint: graphQLEndpoint}
}

type GraphQLServiceImpl struct {
	stockService    ports.InventoryStockService
	graphQLEndpoint string
}

// Defining the Graphql handler
func (g *GraphQLServiceImpl) graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	resolver := &graph.Resolver{StockService: g.stockService}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func (g *GraphQLServiceImpl) playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", g.graphQLEndpoint)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

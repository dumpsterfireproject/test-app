package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cycle-labs/test-app/internal/adapters/inventorydb"
	"github.com/cycle-labs/test-app/internal/domain/ports"
	"github.com/genjidb/genji"
	"github.com/gin-gonic/gin"
)

type Option = func(ctx *gin.Context)

func WithBearerToken(token string) Option {
	return func(c *gin.Context) {
		c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

func WithPathParams(pathParams gin.Params) Option {
	return func(c *gin.Context) {
		c.Params = pathParams
	}
}

func WithQueryParams(queryParams url.Values) Option {
	return func(c *gin.Context) {
		c.Request.URL.RawQuery = queryParams.Encode()
	}
}

func MockJsonGet(c *gin.Context, opts ...Option) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	for _, opt := range opts {
		opt(c)
	}
}

func GetTestGinContext(recorder *httptest.ResponseRecorder) (*gin.Context, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	ctx, engine := gin.CreateTestContext(recorder)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	return ctx, engine
}

var testToken = "AnyTokenIsOk"

func InMemoryStockService(t *testing.T) (*genji.DB, ports.InventoryStockService) {
	db, err := genji.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	repository := inventorydb.NewGengiInventoryStockRepository(db)
	err = repository.SeedData(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	stockService, _ := ports.NewInventoryStock(repository)
	return db, stockService
}

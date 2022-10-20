package api

import (
	"github.com/dumpsterfireproject/test-app/internal/domain/ports"
	"github.com/gin-gonic/gin"
)

// Not a REST API, just a generic JSON over HTTP API
type GinJsonAPIEndpoints struct {
	service GinJsonAPIService
}

func NewGinJsonAPIEndpoints(stockService ports.InventoryStockService) *GinJsonAPIEndpoints {
	api := &GinJsonAPIEndpoints{
		service: NewGinJsonAPIService(stockService),
	}
	return api
}

func (api *GinJsonAPIEndpoints) AddInventoryRoutes(root *gin.RouterGroup) {
	root.GET("items", api.service.ListItemsRequestHander())
	root.GET("items/:id", api.service.GetItemRequestHander())
	root.GET("inventory", api.service.ListInventoryRequestHander())
	root.POST("addInventory", api.service.AddInventoryRequestHander())
	root.DELETE("removeInventory", api.service.DeleteInventoryRequestHander())
}

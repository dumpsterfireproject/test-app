package api

import (
	"net/http"
	"strconv"

	"github.com/cycle-labs/test-app/internal/domain/model"
	"github.com/cycle-labs/test-app/internal/domain/ports"
	"github.com/gin-gonic/gin"
)

type GinJsonAPIService interface {
	ListItemsRequestHander() func(c *gin.Context)
	ListInventoryRequestHander() func(c *gin.Context)
	AddInventoryRequestHander() func(c *gin.Context)
	DeleteInventoryRequestHander() func(c *gin.Context)
}

type GinJsonAPIServiceImpl struct {
	stockService ports.InventoryStockService
}

func NewGinJsonAPIService(stockService ports.InventoryStockService) GinJsonAPIService {
	return &GinJsonAPIServiceImpl{stockService: stockService}
}

func (api *GinJsonAPIServiceImpl) ListItemsRequestHander() func(c *gin.Context) {
	return func(c *gin.Context) {
		items, err := api.stockService.ListItems(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, items)
	}
}

func (api *GinJsonAPIServiceImpl) ListInventoryRequestHander() func(c *gin.Context) {
	return func(c *gin.Context) {
		itemID, itemParameterFound := c.GetQuery("itemID")
		var err error
		var inventory []*model.Inventory
		if itemParameterFound {
			inventory, err = api.stockService.ListInventoryForItem(c.Request.Context(), itemID)
		} else {
			inventory, err = api.stockService.ListInventory(c.Request.Context())
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, inventory)
	}
}

func (api *GinJsonAPIServiceImpl) AddInventoryRequestHander() func(c *gin.Context) {
	return func(c *gin.Context) {
		inventory := &model.Inventory{}
		if err := c.ShouldBindJSON(inventory); err != nil {
			_, err := api.stockService.AddInventory(c.Request.Context(), inventory)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			c.Status(http.StatusCreated)
		}
	}
}

func (api *GinJsonAPIServiceImpl) DeleteInventoryRequestHander() func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: Need returns or such for not founds
		itemID, found := c.GetQuery("itemID")
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "itemID parameter is required"})
		}
		location, found := c.GetQuery("location")
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "location parameter is required"})
		}
		status, found := c.GetQuery("status")
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status parameter is required"})
		}
		quantity, found := c.GetQuery("quantity")
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": "quantity parameter is required"})
		}
		qty, err := strconv.Atoi(quantity)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		_, err = api.stockService.RemoveInventory(c.Request.Context(), itemID, location, status, qty)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Status(http.StatusNoContent)
	}
}

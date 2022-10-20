package graph

import "github.com/dumpsterfireproject/test-app/internal/domain/ports"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	StockService ports.InventoryStockService
}

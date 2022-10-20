package ports

import (
	"context"

	"github.com/dumpsterfireproject/test-app/internal/domain/model"
)

type InventoryStockRepository interface {
	GetItem(ctx context.Context, id string) (*model.Item, error)
	ListItems(ctx context.Context) ([]*model.Item, error)
	ListInventory(ctx context.Context) ([]*model.Inventory, error)
	ListInventoryForItem(ctx context.Context, itemID string) ([]*model.Inventory, error)
	AddInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	RemoveInventory(ctx context.Context, itemID string, location string, status string, quantity int) (*model.Inventory, error)
	SeedData(ctx context.Context) error
}

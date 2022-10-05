package ports

import (
	"context"

	"github.com/cycle-labs/test-app/internal/domain/model"
)

type InventoryStockService interface {
	GetItem(ctx context.Context, id string) (*model.Item, error)
	ListItems(ctx context.Context) ([]*model.Item, error)
	ListInventory(ctx context.Context) ([]*model.Inventory, error)
	ListInventoryForItem(ctx context.Context, itemID string) ([]*model.Inventory, error)
	AddInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	RemoveInventory(ctx context.Context, itenID string, location string, status string, quantity int) (*model.Inventory, error)
}

func NewInventoryStock(repository InventoryStockRepository) (InventoryStockService, error) {
	inventoryStockService := &InventoryStockImpl{repository: repository}
	return inventoryStockService, nil
}

type InventoryStockImpl struct {
	repository InventoryStockRepository
}

func (i *InventoryStockImpl) GetItem(ctx context.Context, id string) (*model.Item, error) {
	return i.repository.GetItem(ctx, id)
}

func (i *InventoryStockImpl) ListItems(ctx context.Context) ([]*model.Item, error) {
	return i.repository.ListItems(ctx)
}

func (i *InventoryStockImpl) ListInventory(ctx context.Context) ([]*model.Inventory, error) {
	return i.repository.ListInventory(ctx)
}

func (i *InventoryStockImpl) ListInventoryForItem(ctx context.Context, itemID string) ([]*model.Inventory, error) {
	return i.repository.ListInventoryForItem(ctx, itemID)
}

func (i *InventoryStockImpl) AddInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	return i.repository.AddInventory(ctx, inventory)
}

func (i *InventoryStockImpl) RemoveInventory(ctx context.Context, itemID string, location string, status string, quantity int) (*model.Inventory, error) {
	return i.repository.RemoveInventory(ctx, itemID, location, status, quantity)
}

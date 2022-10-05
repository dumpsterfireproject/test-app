package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/cycle-labs/test-app/internal/adapters/graph/generated"
	"github.com/cycle-labs/test-app/internal/adapters/graph/model"
	domain "github.com/cycle-labs/test-app/internal/domain/model"
)

// AddInventory is the resolver for the addInventory field.
func (r *mutationResolver) AddInventory(ctx context.Context, input model.NewInventory) (*model.Inventory, error) {
	item, err := r.StockService.GetItem(ctx, input.ItemID)
	if err != nil {
		return nil, err
	}
	newInventory := &domain.Inventory{Item: *item, Location: input.Location, Status: input.Status, Quantity: input.Quantity}
	inv, err := r.StockService.AddInventory(ctx, newInventory)
	if err != nil {
		return nil, err
	}
	return toModelInventory(inv), nil
}

// RemoveInventory is the resolver for the removeInventory field.
func (r *mutationResolver) RemoveInventory(ctx context.Context, input model.RemoveInventory) (*model.Inventory, error) {
	inv, err := r.StockService.RemoveInventory(ctx, input.ItemID, input.Location, input.Status, input.Quantity)
	if err != nil {
		return nil, err
	}
	return toModelInventory(inv), nil
}

// Items is the resolver for the items field.
func (r *queryResolver) Items(ctx context.Context) ([]*model.Item, error) {
	domainItems, err := r.StockService.ListItems(ctx)
	if err != nil {
		return nil, err
	}
	items := []*model.Item{}
	for _, item := range domainItems {
		items = append(items, toModelItem(item))
	}
	return items, nil
}

// ItemInventory is the resolver for the itemInventory field.
func (r *queryResolver) ItemInventory(ctx context.Context, itemID string) ([]*model.Inventory, error) {
	domainInventory, err := r.StockService.ListInventoryForItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	inventory := []*model.Inventory{}
	for _, i := range domainInventory {
		inventory = append(inventory, toModelInventory(i))
	}
	return inventory, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func toModelItem(item *domain.Item) *model.Item {
	return &model.Item{
		ID:          item.ID,
		Sku:         item.SKU,
		Description: item.Description,
	}
}

func toModelInventory(inv *domain.Inventory) *model.Inventory {
	return &model.Inventory{
		ID:       inv.Item.ID,
		Item:     toModelItem(&inv.Item),
		Location: inv.Location,
		Status:   inv.Status,
		Quantity: inv.Quantity,
	}

}

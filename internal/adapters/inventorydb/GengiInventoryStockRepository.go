package inventorydb

import (
	"context"
	"sync"

	// "database/sql"

	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
	"github.com/google/uuid"

	"github.com/cycle-labs/test-app/internal/domain/model"
	"github.com/cycle-labs/test-app/internal/domain/ports"
)

// Return a new InventoryStockRepository that utilizes a genjiDB
func NewGengiInventoryStockRepository(db *genji.DB) ports.InventoryStockRepository {
	repo := &GengiInventoryStockRepository{db: db}
	return repo
}

type GengiInventoryStockRepository struct {
	db *genji.DB
	mu sync.Mutex
}

func (repo *GengiInventoryStockRepository) GetItem(ctx context.Context, id string) (*model.Item, error) {
	item := &model.Item{}
	res, err := repo.db.Query(`select id, sku, description from item where id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	err = res.Iterate(func(d types.Document) error {
		err = document.StructScan(d, item)
		if err != nil {
			return err
		}
		return nil
	})
	return item, nil
}

// List all items in the repository
func (repo *GengiInventoryStockRepository) ListItems(ctx context.Context) ([]*model.Item, error) {
	items := []*model.Item{}
	res, err := repo.db.Query(`select id, sku, description from item order by description`)
	if err != nil {
		return items, err
	}
	defer res.Close()

	err = res.Iterate(func(d types.Document) error {
		item := &model.Item{}
		err = document.StructScan(d, item)
		if err != nil {
			return err
		}
		items = append(items, item)
		return nil
	})
	return items, nil
}

// List all inventory.  TODO: Pagination, maybe filtering
func (repo *GengiInventoryStockRepository) ListInventory(ctx context.Context) ([]*model.Inventory, error) {
	inventory, err := repo.listInventory(ctx, "")
	return inventory, err
}

// List all inventory in the repository for the given item ID
func (repo *GengiInventoryStockRepository) ListInventoryForItem(ctx context.Context, itemID string) ([]*model.Inventory, error) {
	inventory, err := repo.listInventory(ctx, itemID)
	return inventory, err
}

// Add inventory to the repository, incrementing the quantity in the case where some already exists for
// the given item, location, and status
func (repo *GengiInventoryStockRepository) AddInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// upsert
	res, err := repo.db.Query(`SELECT id, quantity FROM inventory WHERE item_id = ? and location = ? and status = ?`,
		inventory.Item.ID, inventory.Location, inventory.Status)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	currentQuantity := 0
	id := ""
	err = res.Iterate(func(d types.Document) error {
		err = document.Scan(d, &id, &currentQuantity)
		return err
	})
	if err != nil {
		return nil, err
	}
	if id == "" {
		id = uuid.NewString()
		err = repo.db.Exec(`INSERT INTO inventory (id, item_id, location, status, quantity) VALUES (?,?,?,?,?)`,
			id, inventory.Item.ID, inventory.Location, inventory.Status, inventory.Quantity)
	} else {
		err = repo.db.Exec(`UPDATE inventory SET quantity = ? WHERE id = ?`, currentQuantity+inventory.Quantity, id)
	}
	if err != nil {
		return nil, err
	}
	return repo.getInventory(ctx, inventory.Item.ID, inventory.Location, inventory.Status)
}

// Remove inventory from the location, with a floor of 0.
// Does not delete the document when the quantity reaches 0.
func (repo *GengiInventoryStockRepository) RemoveInventory(ctx context.Context, itemID string, location string, status string, quantity int) (*model.Inventory, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	// I didn't see anything like a decode or such to handle floor
	res, err := repo.db.Query(`SELECT quantity FROM inventory WHERE item_id = ? and location = ? and status = ?`, itemID, location, status)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	currentQuantity := 0
	err = res.Iterate(func(d types.Document) error {
		err = document.Scan(d, &currentQuantity)
		return err
	})
	newQuantity := 0
	if quantity < currentQuantity {
		newQuantity = currentQuantity - quantity
	}
	err = repo.db.Exec(`UPDATE inventory SET quantity = ? WHERE item_id = ? and location = ? and status = ?`, newQuantity, itemID, location, status)
	if err != nil {
		return nil, err
	}
	return repo.getInventory(ctx, itemID, location, status)
}

// Provide seed data used to start the application
func (repo *GengiInventoryStockRepository) SeedData(ctx context.Context) error {
	err := repo.db.Exec(`
	    CREATE TABLE item (
		    id          TEXT PRIMARY KEY,
			sku         TEXT NOT NULL,
	        description TEXT
	    )
    `)
	if err != nil {
		return err
	}
	err = repo.db.Exec(`
	    CREATE TABLE inventory (
		    id       TEXT PRIMARY KEY,
			item_id  TEXT NOT NULL,
			location TEXT NOT NULL,
			status   TEXT NOT NULL,
	        quantity INT  DEFAULT 0
	    )
    `)
	if err != nil {
		return err
	}

	items := []*model.Item{
		{ID: "8e1af20d-7c39-47e2-a70c-3938bcee2e29", SKU: "60760-3400", Description: "Modafinil"},
		{ID: "7a6b2ef2-1d30-4d1a-a492-fa9d3a361fcc", SKU: "55154-4760", Description: "Glimepiride"},
		{ID: "5fb61987-a6f4-49c5-bc9d-3ed3ff825569", SKU: "62756-1830", Description: "Oxcarbazepine"},
		{ID: "0ccb503a-8944-4c07-b47f-2a1daafe0dee", SKU: "50436-6357", Description: "Warfarin Sodium"},
		{ID: "19540199-bad4-4a34-8383-b42bd762afbe", SKU: "49288-0823", Description: "Acetaminophen, guaifenesin, phenylephrine HCl"},
	}

	for _, item := range items {
		err := repo.db.Exec(`INSERT INTO item VALUES ?`, item)
		if err != nil {
			return err
		}
	}

	inventory := []struct {
		id       string
		itemID   string
		location string
		status   string
		quantity int
	}{
		{"25078cd9-1d4c-47eb-9f47-ff6474df0838", "8e1af20d-7c39-47e2-a70c-3938bcee2e29", "Memphis", "Available", 100},
		{"51b65e8e-b7a2-407a-94d5-ead20b5d12e9", "7a6b2ef2-1d30-4d1a-a492-fa9d3a361fcc", "Memphis", "Available", 75},
		{"0f1f9cee-6f46-4b01-b276-2e35fd093f66", "5fb61987-a6f4-49c5-bc9d-3ed3ff825569", "Memphis", "Available", 18},
		{"f9dafd40-1d16-4c36-8cc7-c0970f80322d", "0ccb503a-8944-4c07-b47f-2a1daafe0dee", "Memphis", "Available", 5000},
		{"f9dafd40-1d16-4c36-8cc7-c0970f80322e", "0ccb503a-8944-4c07-b47f-2a1daafe0dee", "Memphis", "Quarantine", 3300},
	}

	for _, inv := range inventory {
		err = repo.db.Exec(`INSERT INTO inventory (id, item_id, location, status, quantity) VALUES (?,?,?,?,?)`,
			inv.id, inv.itemID, inv.location, inv.status, inv.quantity)
		if err != nil {
			return err
		}
	}

	return err
}

func (repo *GengiInventoryStockRepository) listInventory(ctx context.Context, itemID string) ([]*model.Inventory, error) {
	inventory := []*model.Inventory{}
	item, err := repo.GetItem(ctx, itemID)
	var inventoryRes *genji.Result
	if itemID == "" {
		inventoryRes, err = repo.db.Query(`select id, item_id, location, status, quantity from inventory order by location`)
	} else {
		inventoryRes, err = repo.db.Query(
			`select id, item_id, location, status, quantity from inventory where item_id = ? order by location`,
			itemID)
	}
	defer inventoryRes.Close()
	err = inventoryRes.Iterate(func(d types.Document) error {
		_itemID := ""
		inv := &model.Inventory{}
		err = document.Scan(d, &inv.ID, &_itemID, &inv.Location, &inv.Status, &inv.Quantity)
		if err != nil {
			return err
		}
		inv.Item = *item
		inventory = append(inventory, inv)
		return nil
	})
	if err != nil {
		return inventory, err

	}

	return inventory, nil
}

func (repo *GengiInventoryStockRepository) getInventory(ctx context.Context, itemID string, location string, status string) (*model.Inventory, error) {
	item, err := repo.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}
	inventory := &model.Inventory{Item: *item}
	inventoryRes, err := repo.db.Query(
		`select id, location, status, quantity from inventory where item_id = ? and location = ? and status = ?`,
		itemID, location, status)
	defer inventoryRes.Close()
	err = inventoryRes.Iterate(func(d types.Document) error {
		err = document.Scan(d, &inventory.ID, &inventory.Location, &inventory.Status, &inventory.Quantity)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

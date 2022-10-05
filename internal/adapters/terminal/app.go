package terminal

import (
	"context"

	"github.com/cycle-labs/test-app/internal/domain/model"
	"github.com/cycle-labs/test-app/internal/domain/ports"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	batchSize        = 80            // The number of rows loaded per batch.
	inventoryPageKey = "*inventory*" // The name of the Inventory page.
	modalKey         = "*modal*"     // The name of the modal page.
)

type InventoryApp interface {
	// these can be changed to package level
	ListItems() ([]*model.Item, error)
	ListInventoryForItem(itemID string) ([]*model.Inventory, error)
	AddInventory(inventory *model.Inventory) error
	RemoveInventory(sku string, location string, status string, quantity int) error
	SetFocus(p tview.Primitive) *tview.Application
	PresentModal(p tview.Primitive, width, height int)
	DismissModal()
}

type TerminalApp struct {
	stockService        ports.InventoryStockService
	app                 *tview.Application
	pages               *tview.Pages
	inventoryController *inventoryController
}

func NewTerminalApp(stockService ports.InventoryStockService) *TerminalApp {
	app := tview.NewApplication()
	pages := tview.NewPages()
	termApp := &TerminalApp{
		stockService: stockService,
		app:          app,
		pages:        pages,
	}
	termApp.inventoryController = newInventoryController(termApp)

	flex := termApp.inventoryController.NewInventoryView()

	termApp.pages = tview.NewPages().
		AddPage(inventoryPageKey, flex, true, true)

	termApp.app.SetRoot(termApp.pages, true)
	termApp.setInputCapture()

	return termApp
}

func (t *TerminalApp) Run() error {
	err := t.app.EnableMouse(true).Run()
	return err
}

func (t *TerminalApp) ListItems() ([]*model.Item, error) {
	items, err := t.stockService.ListItems(context.Background())
	return items, err
}

func (t *TerminalApp) ListInventoryForItem(itemID string) ([]*model.Inventory, error) {
	inventory, err := t.stockService.ListInventoryForItem(context.Background(), itemID)
	return inventory, err
}

func (t *TerminalApp) AddInventory(inventory *model.Inventory) error {
	_, err := t.stockService.AddInventory(context.Background(), inventory)
	return err
}

func (t *TerminalApp) RemoveInventory(itemID string, location string, status string, quantity int) error {
	_, err := t.stockService.RemoveInventory(context.Background(), itemID, location, status, quantity)
	return err
}

func (t *TerminalApp) SetFocus(p tview.Primitive) *tview.Application {
	return t.app.SetFocus(p)
}

func (t *TerminalApp) PresentModal(p tview.Primitive, width, height int) {
	modal := createModal(p, width, height)
	t.pages.AddPage(modalKey, modal, true, true)
}

func (t *TerminalApp) DismissModal() {
	t.pages.RemovePage(modalKey)
}

func createModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func (t *TerminalApp) setInputCapture() {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// F1 to quit
		if event.Key() == tcell.KeyF1 {
			t.app.Stop()
		}
		return event
	})
}

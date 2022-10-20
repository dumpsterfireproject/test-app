package terminal

import (
	"fmt"
	"strconv"

	"github.com/dumpsterfireproject/test-app/internal/domain/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type inventoryController struct {
	app        InventoryApp
	itemList   *tview.List
	inventory  *tview.Table
	itemLookup []*model.Item
	modal      tview.Primitive
}

func newInventoryController(app InventoryApp) *inventoryController {
	return &inventoryController{
		app: app,
	}
}

// Create a new tview.Flex element that includes the list of items and inventory table
func (c *inventoryController) NewInventoryView() *tview.Flex {

	// Create the basic objects.
	c.itemList = tview.NewList().ShowSecondaryText(false)
	c.itemList.SetBorder(true).SetTitle("Items")

	c.inventory = tview.NewTable().SetFixed(1, 0)
	c.inventory.SetBorder(true).SetTitle("Inventory")
	c.inventory.SetSelectable(true, false)
	c.inventory.SetDoneFunc(func(key tcell.Key) {
		c.app.SetFocus(c.itemList)
	})

	_items, err := c.app.ListItems()
	if err != nil {
		panic(err)
	}
	c.itemLookup = _items

	for _, item := range _items {
		txt := fmt.Sprintf("%s: %s", item.SKU, item.Description)
		c.itemList.AddItem(txt, "", 0, func() {
			c.app.SetFocus(c.inventory)
		})
		c.itemList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			c.displayInventoryForItem(index)
		})
	}

	// Create the layout.
	flex := tview.NewFlex().
		AddItem(c.itemList, 0, 1, true).
		AddItem(c.inventory, 0, 3, true)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// F2 to add inventory, F3 to remove inventory
		switch event.Key() {
		case tcell.KeyF2:
			c.presentAddInventoryModal()
		case tcell.KeyF3:
			c.presentRemoveInventoryModal()
		}
		return event
	})

	c.displayInventoryForItem(0)

	return flex
}

func (c *inventoryController) displayInventoryForItem(index int) {
	item := c.itemLookup[index]
	c.inventory.Clear()
	c.inventory.SetCell(0, 0, &tview.TableCell{Text: "Location", Align: tview.AlignLeft, Color: tcell.ColorYellow, NotSelectable: true, Expansion: 1})
	c.inventory.SetCell(0, 1, &tview.TableCell{Text: "Status", Align: tview.AlignLeft, Color: tcell.ColorYellow, NotSelectable: true, Expansion: 1})
	c.inventory.SetCell(0, 2, &tview.TableCell{Text: "Quantity", Align: tview.AlignLeft, Color: tcell.ColorYellow, NotSelectable: true, Expansion: 3})

	_inventory, err := c.app.ListInventoryForItem(item.ID)
	if err != nil {
		fmt.Println(err)
	}
	for i, inv := range _inventory {
		quantity := fmt.Sprintf("%d", inv.Quantity)
		selectedStyle := tcell.StyleDefault.Foreground(tview.Styles.PrimitiveBackgroundColor).Background(tview.Styles.PrimaryTextColor)
		c.inventory.SetCell(i+1, 0, &tview.TableCell{Text: inv.Location, Align: tview.AlignLeft}).SetSelectedStyle(selectedStyle)
		c.inventory.SetCell(i+1, 1, &tview.TableCell{Text: inv.Status, Align: tview.AlignLeft}).SetSelectedStyle(selectedStyle)
		c.inventory.SetCell(i+1, 2, &tview.TableCell{Text: quantity, Align: tview.AlignLeft}).SetSelectedStyle(selectedStyle)
	}
}

func (c *inventoryController) presentAddInventoryModal() {
	form := c.createBaseInventoryModal()
	form.AddButton("Save", func() {
		i, _ := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
		if i < 0 {
			return
		}
		item := c.itemLookup[i]
		location := form.GetFormItem(1).(*tview.InputField).GetText()
		status := form.GetFormItem(2).(*tview.InputField).GetText()
		quantity, _ := strconv.Atoi(form.GetFormItem(3).(*tview.InputField).GetText())
		addedInventory := &model.Inventory{
			Item:     *item,
			Location: location,
			Status:   status,
			Quantity: quantity,
		}
		c.app.AddInventory(addedInventory)
		c.dismissModal()
	})
	form.AddButton("Cancel", func() {
		c.dismissModal()
	})

	form.SetBorder(true).SetTitle("Add Inventory").SetTitleAlign(tview.AlignCenter)
	c.app.PresentModal(form, 40, 13)
	c.modal = form
}

func (c *inventoryController) presentRemoveInventoryModal() {
	form := c.createBaseInventoryModal()
	form.AddButton("Save", func() {
		i, _ := form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
		itemID := ""
		if i >= 0 {
			itemID = c.itemLookup[i].ID
		}
		location := form.GetFormItem(1).(*tview.InputField).GetText()
		status := form.GetFormItem(2).(*tview.InputField).GetText()
		quantity, _ := strconv.Atoi(form.GetFormItem(3).(*tview.InputField).GetText())
		c.app.RemoveInventory(itemID, location, status, quantity)
		c.dismissModal()
	})
	form.AddButton("Cancel", func() {
		c.dismissModal()
	})

	form.SetBorder(true).SetTitle("Remove Inventory").SetTitleAlign(tview.AlignCenter)
	c.app.PresentModal(form, 40, 13)
	c.modal = form
}

func (c *inventoryController) createBaseInventoryModal() *tview.Form {
	items := c.itemDescriptions()
	form := tview.NewForm().AddDropDown("Item", items, 0, nil).
		AddInputField("Location", "", 20, nil, nil).
		AddInputField("Status", "", 20, nil, nil).
		AddInputField("Quantity", "", 20, fieldIsNumber, nil)
	return form
}

func fieldIsNumber(textToCheck string, lastChar rune) bool {
	return lastChar >= '0' && lastChar <= '9'
}

func (c *inventoryController) itemDescriptions() []string {
	items := []string{}
	for _, item := range c.itemLookup {
		items = append(items, item.Description)
	}
	return items
}

func (c *inventoryController) dismissModal() {
	c.app.DismissModal()
	c.modal = nil
	// refresh
	i := c.itemList.GetCurrentItem()
	c.displayInventoryForItem(i)
}

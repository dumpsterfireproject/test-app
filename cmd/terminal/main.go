package main

import (
	"context"
	"log"

	"github.com/cycle-labs/test-app/internal/adapters/inventorydb"
	"github.com/cycle-labs/test-app/internal/adapters/terminal"
	"github.com/cycle-labs/test-app/internal/domain/ports"

	"github.com/genjidb/genji"
)

func main() {
	// f, _ := os.OpenFile("testlogfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// log.SetOutput(f)
	// defer f.Close()

	ctx := context.Background()
	db, err := createInMemoryDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repository := inventorydb.NewGengiInventoryStockRepository(db)
	err = repository.SeedData(ctx)
	if err != nil {
		log.Fatal(err)
	}
	stockService, _ := ports.NewInventoryStock(repository)

	theApp := terminal.NewTerminalApp(stockService)
	if err := theApp.Run(); err != nil {
		panic(err)
	}
}

func createInMemoryDatabase() (*genji.DB, error) {
	db, err := genji.Open(":memory:")
	return db, err
}

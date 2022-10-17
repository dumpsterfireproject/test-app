package application

import "github.com/genjidb/genji"

func CreateInMemoryDatabase() (*genji.DB, error) {
	db, err := genji.Open(":memory:")
	return db, err
}

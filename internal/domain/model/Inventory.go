package model

type Inventory struct {
	ID       string `json:"id"`
	Item     Item   `json:"item"`
	Location string `json:"location"`
	Status   string `json:"status"`
	Quantity int    `json:"quantity"`
}

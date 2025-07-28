package entities

type OperationType string

const (
	OperationTypeBuy  OperationType = "buy"
	OperationTypeSell OperationType = "sell"
)

type Operation struct {
	Type     OperationType `json:"operation"`
	UnitCost int           `json:"unit_cost"`
	Quantity int           `json:"quantity"`
}

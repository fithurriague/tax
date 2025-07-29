package entities

type OperationType string

const (
	OperationTypeBuy  OperationType = "buy"
	OperationTypeSell OperationType = "sell"
)

var AllOperationTypes = []OperationType{OperationTypeBuy, OperationTypeSell}

type Operation struct {
	Type     OperationType `json:"operation"`
	UnitCost float64       `json:"unit_cost"`
	Quantity int           `json:"quantity"`
}

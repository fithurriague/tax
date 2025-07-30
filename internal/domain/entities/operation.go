package entities

type OperationType string

const (
	OperationTypeBuy  OperationType = "buy"
	OperationTypeSell OperationType = "sell"
)

var AllOperationTypes = []OperationType{OperationTypeBuy, OperationTypeSell}

type Operation struct {
	Type     OperationType `json:"operation"`
	UnitCost float64       `json:"unit-cost"`
	Quantity int           `json:"quantity"`
}

func (op *Operation) Total() float64 {
	return op.UnitCost * float64(op.Quantity)
}

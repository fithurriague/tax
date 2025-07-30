package entities

import "errors"

type MarketSession struct {
	ActionsBuyed                int
	AccumulatedWeightedUnitCost float64
	AccumulatedLosses           float64
}

func (s *MarketSession) Buy(unitCost float64, quantity int) {
	s.AccumulatedWeightedUnitCost += unitCost * float64(quantity)
	s.ActionsBuyed += quantity
}

func (s *MarketSession) Sell(unitCost float64, quantity int) (profit float64, err error) {
	// Safe to assume it will never happen, but for correctness sake
	if quantity > s.ActionsBuyed {
		return 0, errors.New("not enough actions to sell")
	}

	// Weighted average calculation
	averageUnitCost := s.AccumulatedWeightedUnitCost / float64(s.ActionsBuyed)

	// Accumulate losses
	if unitCost <= averageUnitCost {
		diff := averageUnitCost - unitCost
		s.AccumulatedLosses += diff * float64(quantity)

		return 0, nil
	}

	// Return profit
	diff := unitCost - averageUnitCost
	return diff * float64(quantity), nil
}

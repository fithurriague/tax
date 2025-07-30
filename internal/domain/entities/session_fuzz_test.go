package entities

import (
	"math"
	"testing"
)

func FuzzMarketSessionSell(f *testing.F) {
	f.Add(100.0, 10, 90.0, 5)
	f.Add(0.0, 1, 50.0, 1)
	f.Add(-100.0, -10, -50.0, -5)
	f.Add(0.5, -11, -79.0, 254)
	f.Add(math.MaxFloat64, math.MaxInt, 1.0, 1)
	
	f.Fuzz(func(t *testing.T, buyUnitCost float64, buyQuantity int, sellUnitCost float64, sellQuantity int) {
		session := &MarketSession{}
		session.Buy(buyUnitCost, buyQuantity)
		_, err := session.Sell(sellUnitCost, sellQuantity)
		if err != nil {
			return
		}
	})
}

func FuzzMarketSessionBuySell(f *testing.F) {
	f.Add(100.0, 10, 110.0, 5)
	f.Add(80.0, -90, 17.0, -1)
	f.Add(-50.0, 20, 45.0, 15)
	f.Add(0.0, 0, 100.0, 1)
	f.Add(math.MaxFloat64, math.MaxInt, 1.0, 1)
	
	f.Fuzz(func(t *testing.T, buyUnitCost float64, buyQuantity int, sellUnitCost float64, sellQuantity int) {
		session := &MarketSession{}
		session.Buy(buyUnitCost, buyQuantity)
		_, err := session.Sell(sellUnitCost, sellQuantity)
		if err != nil {
			return
		}
	})
}

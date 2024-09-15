package main

import (
	"math"
)

// TradeMessage represents a trade message received from the WebSocket.
type TradeMessage struct {
	Price string `json:"p"`
}

// GARCHModel represents the GARCH model parameters.
type GARCHModel struct {
	Omega float64
	Alpha float64
	Beta  float64
}

// Predict computes the next volatility prediction using the GARCH model.
func (m *GARCHModel) Predict(data []float64) float64 {
	n := len(data)
	var sigma2 []float64
	sigma2 = append(sigma2, variance(data))

	for i := 1; i < n; i++ {
		newSigma2 := m.Omega + m.Alpha*math.Pow(data[i-1], 2) + m.Beta*sigma2[i-1]
		sigma2 = append(sigma2, newSigma2)
	}

	return math.Sqrt(sigma2[n-1])
}

// OrderSide represents the side of an order (buy or sell).
type OrderSide string

const (
	Buy  OrderSide = "buy"
	Sell OrderSide = "sell"
)

// Order represents a trading order.
type Order struct {
	TickerIndexPut     int       // The index when the order was placed
	TickerIndexExpired int       // The index when the order expires
	TickerIndexFilled  int       // The index when the order was filled
	Price              float64   // The price at which the order triggers
	Side               OrderSide // The side of the order (buy or sell)
}

// Wallet represents a trading wallet with coins and USDTs.
type Wallet struct {
	Coins float64
	USDTs float64
}

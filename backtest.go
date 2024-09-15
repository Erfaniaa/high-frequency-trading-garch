package main

import "log"

func simulateBacktest(prices []float64, orders []Order, wallet Wallet, tradeAmountUSDT float64) []float64 {
	worth := make([]float64, len(prices))
	isOrderAvailable := make([]bool, len(orders))
	for i := range orders {
		isOrderAvailable[i] = true
	}
	for i := 0; i < len(prices); i++ {
		for k := 0; k < NumberOfOrdersPerTick; k++ {
			minValidSellOrderIndex := -1
			maxValidBuyOrderIndex := -1
			for j, order := range orders {
				if !isOrderAvailable[j] {
					continue
				}
				if i >= order.TickerIndexPut && i <= order.TickerIndexExpired {
					switch order.Side {
					case Buy:
						if prices[i] < order.Price && wallet.USDTs >= tradeAmountUSDT {
							if maxValidBuyOrderIndex == -1 || order.Price > orders[maxValidBuyOrderIndex].Price {
								maxValidBuyOrderIndex = j
							}
						}
					case Sell:
						if prices[i] > order.Price && wallet.Coins >= tradeAmountUSDT/order.Price {
							if minValidSellOrderIndex == -1 || order.Price < orders[minValidSellOrderIndex].Price {
								minValidSellOrderIndex = j
							}
						}
					}
				}
			}
			if maxValidBuyOrderIndex != -1 {
				orders[maxValidBuyOrderIndex].TickerIndexFilled = i
				matchedBuyOrder := orders[maxValidBuyOrderIndex]
				wallet.USDTs -= tradeAmountUSDT
				wallet.Coins += tradeAmountUSDT / matchedBuyOrder.Price
				isOrderAvailable[maxValidBuyOrderIndex] = false
				log.Printf("Buy order matched: %+v\n", matchedBuyOrder)
				log.Printf("Wallet: %+v\n", wallet)
			}
			if minValidSellOrderIndex != -1 {
				orders[minValidSellOrderIndex].TickerIndexFilled = i
				matchedSellOrder := orders[minValidSellOrderIndex]
				wallet.USDTs += tradeAmountUSDT
				wallet.Coins -= tradeAmountUSDT / matchedSellOrder.Price
				isOrderAvailable[minValidSellOrderIndex] = false
				log.Printf("Sell order matched: %+v\n", matchedSellOrder)
				log.Printf("Wallet: %+v\n", wallet)
			}
		}
		worth[i] = wallet.Coins*prices[i] + wallet.USDTs
	}
	return worth
}

func calculateStrategyPerformance(prices []float64, orders []Order) (float64, float64, []float64) {
	filledOrdersCount := 0
	for _, order := range orders {
		if order.TickerIndexFilled != -1 {
			filledOrdersCount += 1
		}
	}
	spreadSum := 0.0
	spreadPercents := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		for _, order := range orders {
			if i >= order.TickerIndexPut && i <= order.TickerIndexExpired {
				switch order.Side {
				case Buy:
					spreadSum += (prices[i] - order.Price) / prices[i]
					spreadPercents[i] += 100.0 * (prices[i] - order.Price) / prices[i]
				case Sell:
					spreadSum += (order.Price - prices[i]) / prices[i]
					spreadPercents[i] += 100.0 * (order.Price - prices[i]) / prices[i]
				}
			}
		}
	}
	spreadMeanPercent := 100.0 * spreadSum / float64(len(prices))
	filledOrdersPercent := 100.0 * float64(filledOrdersCount) / float64(len(orders))
	return spreadMeanPercent, filledOrdersPercent, spreadPercents
}

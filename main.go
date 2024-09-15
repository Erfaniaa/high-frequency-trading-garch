package main

import (
	"encoding/json"
	"log"
	"math"
	"net/url"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/websocket"
)

func main() {
	// Set up WebSocket connection
	u := url.URL{Scheme: "wss", Host: "stream.binance.com:9443", Path: "/ws/" + PairName + "@trade"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var prices, lastPrices []float64
	var orders []Order
	garchModel := GARCHModel{Omega: GARCHOmega, Alpha: GARCHAlpha, Beta: GARCHBeta}

	go func() {
		for i := 0; i < MaxIterations; i++ {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var tradeMsg TradeMessage
			if err := json.Unmarshal(message, &tradeMsg); err != nil {
				log.Println("json unmarshal:", err)
				continue
			}

			price, err := strconv.ParseFloat(tradeMsg.Price, 64)
			if err != nil {
				log.Println("price conversion:", err)
				continue
			}
			prices = append(prices, price)

			lastPrices = append(lastPrices, price)
			if len(lastPrices) == GARCHWindowSize+1 {
				lastPrices = lastPrices[1:]
			}

			volatilityPrediction := -1.0
			if len(lastPrices) >= GARCHWindowSize && i%LimitOrderExpiration == 0 {
				garchWindowReturns := make([]float64, len(lastPrices)-1)
				for j := 1; j < len(lastPrices); j++ {
					garchWindowReturns[j-1] = math.Log(lastPrices[j] / lastPrices[j-1])
				}
				garchWindowReturnsAverage := mean(garchWindowReturns)
				volatilityPrediction = garchModel.Predict(garchWindowReturns)
				buyOrder := Order{
					TickerIndexPut:     i + 1,
					TickerIndexExpired: i + LimitOrderExpiration,
					TickerIndexFilled:  -1,
					Price:              prices[i] * math.Exp(garchWindowReturnsAverage-volatilityPrediction*GARCHZScoreCoefficient),
					Side:               Buy,
				}
				sellOrder := Order{
					TickerIndexPut:     i + 1,
					TickerIndexExpired: i + LimitOrderExpiration,
					TickerIndexFilled:  -1,
					Price:              prices[i] * math.Exp(garchWindowReturnsAverage+volatilityPrediction*GARCHZScoreCoefficient),
					Side:               Sell,
				}
				orders = append(orders, buyOrder, sellOrder)
			}

			log.Printf("Iteration Number: %d, Last Price: %f, Volatility Prediction (STDDev): %.12f", i, price, volatilityPrediction)
		}
		interrupt <- os.Interrupt
	}()

	select {
	case <-interrupt:
		log.Println("Received interrupt signal, closing connection...")
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
		}

		wallet := Wallet{USDTs: WalletUSDTs, Coins: WalletCoins}
		tradeAmountUSDT := wallet.USDTs * float64(TradePortionPercent) / 100.0
		worth := simulateBacktest(prices, orders, wallet, tradeAmountUSDT)

		strategyMeanSpreadPercent, strategyFilledOrdersPercent, spreadPercents := calculateStrategyPerformance(prices, orders)

		log.Printf("Strategy Mean Spread Percent: %f", strategyMeanSpreadPercent)
		log.Printf("Strategy Filled Orders Percent: %f", strategyFilledOrdersPercent)
		log.Printf("Wallet Return Percent: %f", 100.0*(worth[len(worth)-1]-worth[0])/worth[0])
		log.Printf("Prices Return Percent: %f", 100.0*(prices[len(prices)-1]-prices[0])/prices[0])

		visualizePricesAndOrders(prices, orders)
		visualizeSpreadPercents(spreadPercents[GARCHWindowSize+1:], GARCHWindowSize+1)
		visualizeWalletWorth(worth)
	}
}

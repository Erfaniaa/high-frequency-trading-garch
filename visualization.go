package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// CustomTicks allows for customized tick marks on plots.
type CustomTicks struct {
	min, max  float64
	precision int
}

func (ct CustomTicks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	nTicks := 10
	spacing := (ct.max - ct.min) / float64(nTicks)
	for i := 0; i <= nTicks; i++ {
		v := ct.min + spacing*float64(i)
		ticks = append(ticks, plot.Tick{
			Value: v,
			Label: fmt.Sprintf("%.*f", ct.precision, v),
		})
	}
	return ticks
}

func visualizePricesAndOrders(prices []float64, orders []Order) {
	p := plot.New()
	p.Title.Text = strings.ToUpper(PairName) + " Prices"
	p.X.Label.Text = "Iteration"
	p.Y.Label.Text = "Price (USDT)"

	var orderPrices []float64
	for _, order := range orders {
		orderPrices = append(orderPrices, order.Price)
	}

	p.Y.Tick.Marker = CustomTicks{min: math.Min(min(prices), min(orderPrices)), max: math.Max(max(prices), max(orderPrices))}

	pts := make(plotter.XYs, len(prices))
	for i, price := range prices {
		pts[i].X = float64(i)
		pts[i].Y = price
	}
	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatalf("Could not create line plot: %v", err)
	}
	p.Add(line)

	for _, order := range orders {
		orderPts := make(plotter.XYs, 2)
		orderPts[0].X = float64(order.TickerIndexPut)
		orderPts[0].Y = order.Price
		orderPts[1].X = float64(order.TickerIndexExpired)
		orderPts[1].Y = order.Price

		orderLine, err := plotter.NewLine(orderPts)
		if err != nil {
			log.Fatalf("Could not create line plot for order: %v", err)
		}

		if order.Side == Buy {
			orderLine.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Green for buy orders
		} else {
			orderLine.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red for sell orders
		}

		if order.TickerIndexFilled == -1 {
			orderLine.Width = vg.Points(1)
		} else {
			orderLine.Width = vg.Points(5)
		}
		p.Add(orderLine)
	}

	if err := p.Save(20*vg.Inch, 10*vg.Inch, "prices.png"); err != nil {
		log.Fatalf("Could not save plot: %v", err)
	}
	log.Println("Plot saved to prices.png")
}

func visualizeSpreadPercents(spreadPercents []float64, skip int) {
	p := plot.New()
	p.Title.Text = "Spread Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Spread Percent"
	p.Y.Tick.Marker = CustomTicks{min: min(spreadPercents), max: max(spreadPercents), precision: 6}
	points := make(plotter.XYs, len(spreadPercents))
	for i := range points {
		points[i].X = float64(skip + i)
		points[i].Y = spreadPercents[i]
	}
	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	p.Add(line)
	if err := p.Save(20*vg.Inch, 10*vg.Inch, "spread_percents.png"); err != nil {
		panic(err)
	}
	log.Println("Spread percents plot saved to spread_percents.png")
}

func visualizeWalletWorth(worth []float64) {
	p := plot.New()
	p.Title.Text = "Wallet Worth Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Wallet Worth (USDTs)"
	p.Y.Tick.Marker = CustomTicks{min: min(worth), max: max(worth), precision: 6}
	points := make(plotter.XYs, len(worth))
	for i := range points {
		points[i].X = float64(i)
		points[i].Y = worth[i]
	}
	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}
	p.Add(line)
	if err := p.Save(20*vg.Inch, 10*vg.Inch, "wallet_worth.png"); err != nil {
		panic(err)
	}
	log.Println("Wallet worth plot saved to wallet_worth.png")
}

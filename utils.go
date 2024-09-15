package main

import "math"

func mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

func variance(data []float64) float64 {
	meanValue := mean(data)
	sum := 0.0
	for _, value := range data {
		sum += math.Pow(value-meanValue, 2)
	}
	return sum / float64(len(data))
}

func varianceEstimate(data []float64) float64 {
	meanValue := mean(data)
	sum := 0.0
	for _, value := range data {
		sum += math.Pow(value-meanValue, 2)
	}
	return sum / float64(len(data)-1)
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minValue := values[0]
	for _, v := range values {
		if v < minValue {
			minValue = v
		}
	}
	return minValue
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxValue := values[0]
	for _, v := range values {
		if v > maxValue {
			maxValue = v
		}
	}
	return maxValue
}

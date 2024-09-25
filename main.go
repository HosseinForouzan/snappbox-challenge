package main
import (
	"time"
	"sync"
)

type Point struct {
	DeliveryID string
	Lat float64
	Lng float64
	Timestamp time.Time
}

type Segment struct {
	P1 Point
	P2 Point
	Speed float64
	Distance float64
	Duration float64
}

type FareCalculator struct {
	mu sync.Mutex
	fares map[string]float64
	waitGroup sync.WaitGroup
}


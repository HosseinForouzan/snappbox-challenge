package main
import (
	"time"
	"sync"
	"math"
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

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 //Radius of the Earth
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	a := (math.Sin(dLat / 2) * math.Sin(dLat / 2)) + 
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c

	return d
} 

func getState(speed float64) string {
	if speed > 10 {
		return "MOVING"
	}
	return "IDLE"
}

func getRate(segment Segment) float64 {
	state := getState(segment.Speed)
	if state == "MOVING" {
		hour := segment.P1.Timestamp.Hour()
		if hour >= 5 && hour < 24 {
			return 0.74 * segment.Distance
		} else {
			return 1.30 * segment.Distance
		}
	} else {
		return 11.90 * segment.Duration
	}
}


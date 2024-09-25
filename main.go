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

func processDelivery(deliveryID string, points []Point, fareCalculator *FareCalculator) {
	defer fareCalculator.waitGroup.Done()

	var totalFare float64 = 1.30
	var validPoints []Point
	
	for i := 0; i < len(points) - 1; i++ {
		p1 := points[i]
		p2 := points[i+1]

		duration := p2.Timestamp.Sub(p1.Timestamp).Hours()
		if duration <= 0 {
			continue
		}

		distance := haversine(p1.Lat, p1.Lng, p2.Lat, p2.Lng)
		speed := distance / duration

		if speed > 100 {
			continue
		}
		validPoints = append(validPoints, p1)
	}

	//Validate last point
	if len(points) > 0 {
		validPoints = append(validPoints, points[len(points) - 1])
	}

	for i := 0; i < len(validPoints) - 1; i++ {
		p1 := validPoints[i]
		p2 := validPoints[i + 1]

		duration := p2.Timestamp.Sub(p1.Timestamp).Hours()
		if duration <= 0 {
			continue
		}

		distance := haversine(p1.Lat, p1.Lng, p2.Lat, p2.Lng)
		speed := distance / duration

		segment := Segment {
			P1: p1,
			P2: p2,
			Speed: speed,
			Distance: distance,
			Duration: duration,
		}
		
		rate := getRate(segment)
		totalFare += rate
	}

	if totalFare < 3.47 {
		totalFare = 3.47
	}

	fareCalculator.mu.Lock()
	fareCalculator.fares[deliveryID] = totalFare
	fareCalculator.mu.Unlock()
}


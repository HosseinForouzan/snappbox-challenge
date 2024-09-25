package main

import (
    "math"
    "testing"
    "time"
)

func TestHaversine(t *testing.T) {
    lat1 := 36.12
    lon1 := -86.67
    lat2 := 33.94
    lon2 := -118.40

    distance := haversine(lat1, lon1, lat2, lon2)
    expected := 2887.26 

    if math.Abs(distance-expected) > 1 {
        t.Errorf("We expected %.2f but got %.2f", distance, expected)
    }
}

func TestGetState(t *testing.T) {
    if state := getState(15); state != "MOVING" {
        t.Errorf("We Expected MOVING but got %s", state)
    }
    if state := getState(5); state != "IDLE" {
        t.Errorf("We expected IDLE but got %s", state)
    }
}

func TestGetRate(t *testing.T) {
    p1 := Point{
        Lat:       35.0,
        Lng:       51.0,
        Timestamp: time.Date(2023, time.October, 1, 6, 0, 0, 0, time.UTC),
    }
    p2 := Point{
        Lat:       35.1,
        Lng:       51.1,
        Timestamp: time.Date(2023, time.October, 1, 6, 10, 0, 0, time.UTC),
    }

    distance := haversine(p1.Lat, p1.Lng, p2.Lat, p2.Lng)
    duration := p2.Timestamp.Sub(p1.Timestamp).Hours()
    speed := distance / duration

    segment := Segment{
        P1:       p1,
        P2:       p2,
        Speed:    speed,
        Distance: distance,
        Duration: duration,
    }

    rate := getRate(segment)
    expectedRate := 0.74 * distance 

    if math.Abs(rate-expectedRate) > 0.01 {
        t.Errorf("We expected %.2f but got %.2f", rate, expectedRate)
    }

    segment.Speed = 5 
    segment.Duration = 0.1667

    rate = getRate(segment)
    expectedRate = 11.90 * segment.Duration

    if math.Abs(rate-expectedRate) > 0.01 {
        t.Errorf("We expected %.2f but got %.2f", rate, expectedRate)
    }
}

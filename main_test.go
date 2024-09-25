package main

import (
    "math"
    "testing"
    "time"
	"os"
	"strings"
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

func TestEndToEnd(t *testing.T) {
    // Create sample input file
    inputFileName := "test_input.csv"
    outputFileName := "test_output.csv"

    inputData := `1,35.0,51.0,1696156800
1,35.1,51.1,1696157400
1,35.2,51.2,1696158000
2,36.0,52.0,1696158600
2,36.1,52.1,1696159200`

    err := os.WriteFile(inputFileName, []byte(inputData), 0644)
    if err != nil {
        t.Fatalf("Error writing input file: %v", err)
    }
    defer os.Remove(inputFileName)
    // defer os.Remove(outputFileName)

    // Run fare calculation
    err = calculateFares(inputFileName, outputFileName)
    if err != nil {
        t.Fatalf("Error executing CalculateFares: %v", err)
    }

    // Read and verify output file
    outputData, err := os.ReadFile(outputFileName)
    if err != nil {
        t.Fatalf("Error reading output file: %v", err)
    }

    outputLines := strings.Split(strings.TrimSpace(string(outputData)), "\n")
    if len(outputLines) != 2 {
        t.Fatalf("Unexpected number of output lines: %d", len(outputLines))
    }

    // Create a map to check results
    fares := make(map[string]string)
    for _, line := range outputLines {
        fields := strings.Split(line, ",")
        if len(fields) != 2 {
            t.Fatalf("Invalid output line format: %s", line)
        }
        fares[fields[0]] = fields[1]
    }

    // Check fare for delivery 1
    if fare, ok := fares["1"]; ok {
        expectedFare := "3.47" // Minimum fare
        if fare != expectedFare {
            t.Errorf("Fare for delivery 1 should be %s, but got %s", expectedFare, fare)
        }
    } else {
        t.Errorf("Delivery 1 not found in output")
    }

    // Check fare for delivery 2
    if fare, ok := fares["2"]; ok {
        expectedFare := "3.47" // Minimum fare
        if fare != expectedFare {
            t.Errorf("Fare for delivery 2 should be %s, but got %s", expectedFare, fare)
        }
    } else {
        t.Errorf("Delivery 2 not found in output")
    }
}
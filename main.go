package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "math"
    "os"
    "strconv"
    "strings"
    "sync"
    "time"
)

type Point struct {
    DeliveryID string
    Lat        float64
    Lng        float64
    Timestamp  time.Time
}

type Segment struct {
    P1       Point
    P2       Point
    Speed    float64 // km/h
    Distance float64 // km
    Duration float64 // hours
}

type FareCalculator struct {
    mu        sync.Mutex
    fares     map[string]float64
    waitGroup sync.WaitGroup
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // Earth radius in km
    dLat := (lat2 - lat1) * math.Pi / 180.0
    dLon := (lon2 - lon1) * math.Pi / 180.0

    lat1 = lat1 * math.Pi / 180.0
    lat2 = lat2 * math.Pi / 180.0

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
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
   
    var totalFare float64 = 1.30 // Initial flag amount
    var validPoints []Point

    // Filter invalid points
    for i := 0; i < len(points)-1; i++ {
        p1 := points[i]
        p2 := points[i+1]

        duration := p2.Timestamp.Sub(p1.Timestamp).Hours()
        if duration <= 0 {
            continue // Invalid duration
        }

        distance := haversine(p1.Lat, p1.Lng, p2.Lat, p2.Lng)
        speed := distance / duration

        if speed > 100 {
            // p2 is invalid and will be removed
            continue
        }

        validPoints = append(validPoints, p1)
    }

    // Add the last point if valid
    if len(points) > 0 {
        validPoints = append(validPoints, points[len(points)-1])
    }

    // Calculate fare using valid points
    for i := 0; i < len(validPoints)-1; i++ {
        p1 := validPoints[i]
        p2 := validPoints[i+1]

        duration := p2.Timestamp.Sub(p1.Timestamp).Hours()
        if duration <= 0 {
            continue
        }

        distance := haversine(p1.Lat, p1.Lng, p2.Lat, p2.Lng)
        speed := distance / duration

        segment := Segment{
            P1:       p1,
            P2:       p2,
            Speed:    speed,
            Distance: distance,
            Duration: duration,
        }

        rate := getRate(segment)
        totalFare += rate
    }

    // Apply minimum fare
    if totalFare < 3.47 {
        totalFare = 3.47
    }

    // Store the calculated fare
    fareCalculator.mu.Lock()
    fareCalculator.fares[deliveryID] = totalFare
    fareCalculator.mu.Unlock()
}

func calculateFares(inputFileName string, outputFileName string) error {
    inputFile, err := os.Open(inputFileName)
    if err != nil {
        return fmt.Errorf("Error in opening input file %v", err)
    }
    defer inputFile.Close()

    scanner := bufio.NewScanner(inputFile)
    fareCalculator := FareCalculator{
        fares: make(map[string]float64),
    }

    var currentDeliveryID string
    var points []Point

    semaphore := make(chan struct{}, 10) // limit goroutines to 10

    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Split(line, ",")

        if len(fields) != 4 {
            continue 
        }

        deliveryID := fields[0]
        lat, err1 := strconv.ParseFloat(fields[1], 64)
        lng, err2 := strconv.ParseFloat(fields[2], 64)
        timestampInt, err3 := strconv.ParseInt(fields[3], 10, 64)

        if err1 != nil || err2 != nil || err3 != nil {
            continue 
        }

        timestamp := time.Unix(timestampInt, 0)

        point := Point{
            DeliveryID: deliveryID,
            Lat:        lat,
            Lng:        lng,
            Timestamp:  timestamp,
        }

        if currentDeliveryID == "" {
            currentDeliveryID = deliveryID
        }

        if deliveryID != currentDeliveryID {
            semaphore <- struct{}{}
            fareCalculator.waitGroup.Add(1)
            go func(deliveryID string, points []Point) {
                defer fareCalculator.waitGroup.Done()
                defer func() { <-semaphore }()
                processDelivery(deliveryID, points, &fareCalculator)
            }(currentDeliveryID, points)

            currentDeliveryID = deliveryID
            points = []Point{point}
        } else {
            points = append(points, point)
        }
    }

    if len(points) > 0 {
        semaphore <- struct{}{}
        fareCalculator.waitGroup.Add(1)
        go func(deliveryID string, points []Point) {
            defer fareCalculator.waitGroup.Done()
            defer func() { <-semaphore }()
            processDelivery(deliveryID, points, &fareCalculator)
        }(currentDeliveryID, points)
    }

    fareCalculator.waitGroup.Wait()

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("Error in reading file %v", err)
    }

    outputFile, err := os.Create(outputFileName)
    if err != nil {
        return fmt.Errorf("Error in creating file %v", err)
    }
    defer outputFile.Close()

    writer := csv.NewWriter(outputFile)

    for deliveryID, fare := range fareCalculator.fares {
        record := []string{deliveryID, fmt.Sprintf("%.2f", fare)}
        if err := writer.Write(record); err != nil {
            return fmt.Errorf("Error in writing record to the file %v", err)
        }
    }

    writer.Flush()
    if err := writer.Error(); err != nil {
        return fmt.Errorf("Error of buffer %v", err)
    }

    return nil
}

func main() {
    err := calculateFares("sample_data.csv", "output.csv")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println("calculating fares has done.")
}

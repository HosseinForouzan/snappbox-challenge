# Delivery Fare Estimation Script

This project is a Go program designed to estimate delivery fares based on GPS data for couriers. It processes input data containing location points for deliveries, filters out invalid data, calculates fares according to specified rules, and outputs the fare estimates for each delivery.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Usage](#usage)
- [Implementation Details](#implementation-details)
- [Testing](#testing)


## Overview

Couriers perform thousands of deliveries daily, and it's crucial to ensure accurate fare calculations and detect any discrepancies. This program automates the fare estimation process by:

- Filtering out invalid GPS data (e.g., unrealistic speeds).
- Calculating delivery fares based on distance, time, and predefined rates.
- Handling large datasets efficiently using concurrency features in Go.

## Features

- **Data Filtering**: Removes GPS points where the calculated speed between two points exceeds 100 km/h.
- **Fare Calculation**:
  - Charges a standard flag amount of **1.30** at the start of each delivery.
  - Applies different rates based on the courier's state (moving or idle) and time of day.
  - Ensures a minimum delivery fare of **3.47**.
- **Concurrency**: Processes multiple deliveries concurrently, utilizing Goroutines and Channels for high performance.
- **Scalability**: Capable of processing large datasets (several gigabytes) efficiently.
- **Error Handling**: Manages errors gracefully and logs informative messages.



## Usage

### 1. Prepare the Input Data

- Place your input CSV file named `input.csv` in the project directory.
- The CSV file should have the following format:

  ```
  id_delivery,lat,lng,timestamp
  ```

  - `id_delivery`: Delivery ID (string).
  - `lat`: Latitude (float64).
  - `lng`: Longitude (float64).
  - `timestamp`: Unix timestamp in seconds (int64).

### 2. Build the Program

```bash
go build -o fare_estimato
```

### 3. Run the Program

```bash
./snapp
```

- The program will read `input.csv` and generate `output.csv` containing fare estimates.

### 4. Check the Output

- The `output.csv` file will have the following format:

  ```
  id_delivery,fare_estimate
  ```

  - `id_delivery`: Delivery ID.
  - `fare_estimate`: Calculated fare in euros, rounded to two decimal places.

## Implementation Details

### Data Structures

- **Point**: Represents a GPS point with latitude, longitude, timestamp, and delivery ID.
- **Segment**: Represents a segment between two consecutive points, including speed, distance, and duration.
- **FareCalculator**: Manages fare calculations and concurrency control using a mutex and wait group.

### Fare Calculation Rules

#### States and Rates

| **State**           | **Applicable When**             | **Fare Amount**               |
|---------------------|---------------------------------|-------------------------------|
| **MOVING (U > 10 km/h)** | Time of day (05:00 to 00:00)    | €0.74 per km                  |
|                     | Time of day (00:00 to 05:00)    | €1.30 per km                  |
| **IDLE (U ≤ 10 km/h)**  | Always                          | €11.90 per hour of idle time   |

- **Flag Amount**: A standard charge of **€1.30** at the start of each delivery.
- **Minimum Fare**: The total fare for each delivery should be at least **€3.47**.

### Data Filtering

- Calculates speed `U` between two consecutive points.
- If `U > 100 km/h`, the second point is considered invalid and removed.

### Concurrency

- Uses Goroutines to process deliveries concurrently.
- Limits the number of concurrent Goroutines using a semaphore to prevent resource exhaustion.
- Synchronizes access to shared data using mutexes.

### Distance Calculation

- Utilizes the **Haversine formula** to calculate the distance between two GPS coordinates.
- Accounts for Earth's curvature for accurate distance measurement.

```go
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // Earth's radius in km
    // Convert degrees to radians
    dLat := (lat2 - lat1) * math.Pi / 180.0
    dLon := (lon2 - lon1) * math.Pi / 180.0
    lat1 = lat1 * math.Pi / 180.0
    lat2 = lat2 * math.Pi / 180.0
    // Apply Haversine formula
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}
```

## Testing

### Unit Tests

- Located in `fare_estimator_test.go`.
- Test individual functions like `haversine`, `getState`, and `getRate`.
- Run unit tests using:

  ```bash
  go test -v -run TestFunctionName
  ```

### End-to-End Test

- Tests the entire flow from reading input to generating output.
- Uses a sample input and checks if the output matches expected fare estimates.
- Run the end-to-end test using:

  ```bash
  go test -v -run TestEndToEnd
  ```

### Running All Tests

```bash
go test -v
```



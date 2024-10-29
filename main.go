package main

import (
    "fmt"
    "log"
    "net/http"
	"encoding/json"
	"io"
)

// Response struct of the Weather Service
type WeatherResponse struct {
	Forecast string `json:"forcast"`
	TemperatureFeel string `json:"temperature_feel"`
}

// Nested struct to get the nested JSON response
type NWSPoint struct {
    Properties struct {
        GridId string `json:"gridId"`
        GridX  int    `json:"gridX"`
        GridY  int    `json:"gridY"`
    } `json:"properties"`
}

// Struct to take the response from the forcast url
type NWSForecast struct {
    Properties struct {
        Periods []struct {
            Temperature     int    `json:"temperature"`
            ShortForecast  string `json:"shortForecast"`
        } `json:"periods"`
    } `json:"properties"`
}


func main() {
    http.HandleFunc("/weather", getWeatherHandler)
    
    port := ":8080"
    fmt.Printf("Server starting on port %s...\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

func getWeatherHandler(w http.ResponseWriter, r *http.Request) {

	latitude := r.URL.Query().Get("lat")
    longitude := r.URL.Query().Get("lon")

    if latitude == "" {
        http.Error(w, "Missing Latitude parameters", http.StatusBadRequest)
        return
    }
	if longitude == "" {
        http.Error(w, "Missing Longitude parameters", http.StatusBadRequest)
        return
    }

	point, err := fetchNWSPoint(latitude, longitude)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	forecast, err := fetchForecast(point)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	response, err := getNowForecast(forecast)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	// Send JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)

}

func categorizeTemperature(temp int) string {
    switch {
    case temp >= 83:
        return "hot"
    case temp <= 63:
        return "cold"
    default:
        return "moderate"
    }
}

func fetchNWSPoint(latitude, longitude string) (*NWSPoint, error) {
    pointsURL := fmt.Sprintf("https://api.weather.gov/points/%s,%s", latitude, longitude)
    
    pointsResp, err := http.Get(pointsURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch grid coordinates: %w", err)
    }
    defer pointsResp.Body.Close()

    // Check for non-200 status code
    if pointsResp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(pointsResp.Body)
        return nil, fmt.Errorf("NWS API returned status %d: %s", 
            pointsResp.StatusCode, string(body))
    }

    var point NWSPoint
    if err := json.NewDecoder(pointsResp.Body).Decode(&point); err != nil {
        return nil, fmt.Errorf("failed to parse grid coordinates response: %w", err)
    }

    return &point, nil
}

func fetchForecast(point *NWSPoint) (*NWSForecast, error) {
    forecastURL := fmt.Sprintf("https://api.weather.gov/gridpoints/%s/%d,%d/forecast",
        point.Properties.GridId,
        point.Properties.GridX,
        point.Properties.GridY)
    
    forecastResp, err := http.Get(forecastURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch forecast: %w", err)
    }
    defer forecastResp.Body.Close()

    // Check for non-200 status code
    if forecastResp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(forecastResp.Body)
        return nil, fmt.Errorf("forecast API returned status %d: %s", 
            forecastResp.StatusCode, string(body))
    }

    var forecast NWSForecast
    if err := json.NewDecoder(forecastResp.Body).Decode(&forecast); err != nil {
        return nil, fmt.Errorf("failed to parse forecast response: %w", err)
    }

    // Validate forecast data
    if len(forecast.Properties.Periods) == 0 {
        return nil, fmt.Errorf("no forecast periods available")
    }

    return &forecast, nil
}


func getNowForecast(forecast *NWSForecast) (*WeatherResponse, error) {
    if forecast == nil {
        return nil, fmt.Errorf("forecast data is nil")
    }
    
    if len(forecast.Properties.Periods) == 0 {
        return nil, fmt.Errorf("no forecast periods available")
    }

    today := forecast.Properties.Periods[0]
    return &WeatherResponse{
        Forecast:    today.ShortForecast,
        TemperatureFeel: categorizeTemperature(today.Temperature),
    }, nil
}

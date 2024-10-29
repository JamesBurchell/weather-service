# Project Submission - Weather Service Assignment

## Write an HTTP server that serves the current weather. 

Your server should expose an endpoint that:
    
1. Accepts latitude and longitude coordinates
2. Returns the short forecast for that area for Today (“Partly Cloudy” etc)
3. Returns a characterization of whether the temperature is “hot”, “cold”, or “moderate” (use your discretion on mapping temperatures to each type)
4. Use the National Weather Service API Web Service as a data source. 

The purpose of this exercise is to provide a sample of your work that we can discuss together in the Technical Interview.
- We respect your time. Spend as long as you need, but we intend it to take around an hour.
- We do not expect a production-ready service, but you might want to comment on your shortcuts.
- The submitted project should build and have brief instructions so we can verify that it works.
- You may write in whatever language or stack you're most comfortable in, but the technical interviewers are most familiar with Typelevel Scala.

## Implementation

- I will use Go as that is a request lang
- Request will be a GET request with params
- Response will be JSON
- I will keep it very simple
- I will forgo testing
- I will aim to complete in about an hour

API:
GET request:
```
http://localhost:8080/weather?lat=<LAT>&lon=<LON>
```
Response:
```
{
    "forecast": "Partly Cloudy",
    "temperature_feel": "moderate"
}
```

### Test
Terminal 1
```
go run .
```

Terminal 2
```
curl "http://localhost:8080/weather?lat=36.174465&lon=-86.767960"
```

Expected response example:
```
{"forcast":"Sunny","temperature_feel":"moderate"}
```
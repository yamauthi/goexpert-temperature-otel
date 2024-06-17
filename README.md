# GoExpert - Temperature by Postal code - With Open telemetry and Zipkin

### Challenge description:
#### Objective
Develop a system in Go that receives a ZIP code, identifies the city and returns the current weather (temperature in degrees celsius, fahrenheit and kelvin) along with the city. This system should implement OTEL(Open Telemetry) and Zipkin.

Based on the known scenario "Temperature system by ZIP code" referred to as Service B, a new project will be included, referred to as Service A.

#### Service A - Input Service
- The system should receive a 8-digits input via POST, using the schema: `{ "cep": "31275000" }`
- The system should validate whether the input is valid (contains 8 digits) and is a STRING
- If valid, it will be forwarded to Service B via HTTP
- If not valid, it should return:
- HTTP Code: 422
- Message: invalid zipcode

#### Service B - Temperature by postal code Service
Receive a brazilian postal code (usually called CEP) and return the temperature on that location in Celsius, Fahrenheit and Kelvin.
- Postal code to be valid must be only numbers, 8 digits
- If a postal code is invalid should return status code 422 with error "invalid postal code"
- If a postal code is valid, but not found shold return status 404 with error "can not find postal code"
- You can get the location using an address API (like ViaCepAPI https://viacep.com.br/ )
- You can get the temperature using a weather API (like WeatherAPI https://www.weatherapi.com/)

### Open Telemetry and Zipkin
Upon implementation of the services, add the implementation of OTEL + Zipkin:
- Implement distributed tracing between Service A - Service B
- Use span to measure the response time of the ZIP code search service and temperature search service

### Requirements
- Valid API Key for WeatherAPI

### How to run locally with docker-compose
- Clone the repository
- Open terminal in project folder
- You can create a .env file from .env.example files inside cmd/input-service for Postal Code Input Service and cmd/temp-postal-code for Temperature by Postal code Service
- You can also change directly into docker-compose.yaml file
- Place your weatherAPI apiKey in cmd/temp-postal-code/.env file or change it on temperature_service definition.
- If any variable is defined on docker-compose.yaml it will overwrite the .env variables 
- Run `docker-compose up --build` and after `docker-compose up -d`
- Default port will be 8080
- Tracing can be checked on Zipkin running on port 9411

### Tests
- You can execute a POST call locally on http://localhost:8080/temperature following the structure mentioned at the beginning

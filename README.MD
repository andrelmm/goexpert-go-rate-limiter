## Rate Limiter

This project implements a rate limiter middleware for HTTP requests using Redis as storage. It allows you to limit the number of requests or actions allowed within a certain time frame for specific clients based on their IP address or API key.

### Features
**IP Address**: The rate limiter restricts the number of requests received from a single IP address within a defined time interval.

**Token**: The rate limiter can also limit requests based on a unique access token, allowing different expiration time limits for different tokens. The token must be provided in the header in the following format:

```
API_KEY: <TOKEN>
```

### Configuration

Before running the application, make sure you have a .env file in the root directory containing the following environment variables:

- RATE_LIMIT_DURATION: Duration for which the rate limit applies (e.g., "1m", "1h").
- RATE_LIMIT_TOKEN: Maximum number of requests allowed within the rate limit duration for API keys.
- RATE_LIMIT_IP: Maximum number of requests allowed within the rate limit duration for IP addresses.
- BLOCK_DURATION: Duration for which a client is blocked after reaching the rate limit (e.g., "5m").

### Example .env file:
```
RATE_LIMIT_IP=5
RATE_LIMIT_TOKEN=10
BLOCK_DURATION=30s
RATE_LIMIT_DURATION=5s
```

### Usage

Run the application using the following command:

```
go run main.go
```

By default, the application listens on port 8080. You can modify the port in the main.go file if needed.

### API
##### Middleware
The rate limiter middleware is applied to all routes by default. It checks the rate limit for each incoming request based on the client's IP address or API key.

### Storage
The storage package provides an interface for interacting with Redis. It includes methods for getting, setting, and manipulating data stored in Redis.

### Running the Application Locally using Docker

To run the application locally using Docker, execute the following commands:

```
docker-compose build
docker-compose up -d
```
To test the application, make a GET request to the following endpoint:

```http request
http://localhost:8080/
``` 
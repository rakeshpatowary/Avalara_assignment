# Avalara_assignment
A very basic URL shortener implemented in GO, using the Gorilla Mux router.

## Overview
This URL shortener provides a RESTful API endpoint to shorten any valid URLs. It generates a short key for each valid URL and allows user to access the original URL by providing the short key.

## Endpoints

### Shorten URL

- **Endpoint:** `PUT /shorten`
- **Request:**
  - Method: `PUT`
  - Body: JSON object with a `destination` field representing the original URL.

```json
{
  "destination": "https://www.google.com"
}

Response:
    - Status Code: `200 OK`
    - Body: JSON object with a short_url field representing the shortened URL

```json
{
  "short_url": "http://avalara-domain.com/JhLLu"
}

### Redirect to Original URL

- **Endpoint:** `GET /{shortKey}`
- **Request:**
  - Method: `GET`
  - Path: The short key generated during URL shortening

Response:
    - Redirect to the original URL if the short key is valid. Otherwise, a 404 Not Found response.

## Usage
1. Clone the repository
2. cd url-shortener_api
3. go build
   ./url-shortener_api
4. Access the API using tools like POSTMAN or CURL
    e.g. If you want to use CURL command
    # Shorten URL
    curl -X PUT -d '{"destination": "https://www.google.com"}' http://localhost:9090/shortURL

    # Redirect to Original URL
    curl http://localhost:9090/abcde

    Note: You may get err if you are running the above command in Windows machine.
    For windowa machine:
    curl -X PUT -d "{\"destination\": \"https://www.google.com\"}" http://localhost:9090/shortURL

## Dependencies
    - Gorilla Mux: A powerful URL router and dispatcher

## License


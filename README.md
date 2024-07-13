# Snap-Link

Snap-Link is a powerful and easy-to-use URL shortener written in Go. This project provides a simple web interface and API for generating and managing shortened URLs.

## Features

- **URL Shortening:** Easily create shortened URLs.
- **API Access:** Use the provided API for programmatic access to URL shortening.
- **Custom Aliases:** Create custom aliases for your shortened URLs.
- **Statistics:** Track the number of clicks on your shortened URLs.
- **Simple Web Interface:** User-friendly web interface for managing links.

## Project Status

This project is still under development and is currently in its early version. New features and improvements will be added over time.

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Docker (optional, for containerized deployment)
- SQLite database
- Configuration file `config.yaml`

### Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/velvetriddles/snap-link.git
    ```

2. **Navigate to the project directory:**

    ```bash
    cd snap-link
    ```

3. **Install dependencies:**

    ```bash
    go mod download
    ```

4. **Create the `config.yaml` configuration file:**

    ```yaml
    env: local
    storage_path: "snap-link.db"
    http_server:
      address: ":8080"
      timeout: 4s
      idle_timeout: 60s
      user: "your_username"
      password: "your_password"
    ```

5. **Run the application:**

    ```bash
    CONFIG_PATH=config.yaml go run main.go
    ```

### Using Docker

1. **Create the `config.yaml` configuration file:**

    ```yaml
    env: local
    storage_path: "snap-link.db"
    http_server:
      address: ":8080"
      timeout: 4s
      idle_timeout: 60s
      user: "your_username"
      password: "your_password"
    ```

2. **Build the Docker image:**

    ```bash
    docker build -t snap-link .
    ```

3. **Run the container:**

    ```bash
    docker run -d -p 8080:8080 -v $(pwd)/config.yaml:/app/config.yaml --name snap-link snap-link
    ```

## API Documentation

### Create a Shortened URL

- **Endpoint:** `/url/`
- **Method:** `POST`
- **Request Parameters:**
  - `url` (string): The URL to shorten
  - `alias` (string, optional): Custom alias for the shortened URL

- **Example Request:**

    ```json
    {
      "url": "https://example.com",
      "alias": "example"
    }
    ```

- **Example Response:**

    ```json
    {
      "short_url": "http://localhost:8080/example"
    }
    ```

### Redirect to Original URL

- **Endpoint:** `/{alias}`
- **Method:** `GET`

Redirects the user to the original URL.

## Contributing

We welcome contributions to the project! Please follow standard fork and pull request processes, and provide detailed descriptions of your changes.

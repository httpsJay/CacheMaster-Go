# Pokemon Cache API

Welcome to the project's Pokemon Cache API! In this project, a cache for Pokemon data is implemented and made available via a REST API. The least frequently utilised items are removed from the cache when its maximum capacity is reached.

## Features

- Store and retrieve Pokemon data by name or ID
- Delete Pokemon data by ID
- Handle concurrent requests gracefully
- Error handling and input validation

## Prerequisites

- Golang
- Docker (optional, for containerization)

## Installation

1. **Install dependencies:**

    ```sh
    make install-deps
    ```

## Running the Application

### Using Makefile

1. **Build the application:**

    ```sh
    make build
    ```

2. **Run the application:**

    ```sh
    make run
    ```

3. **Run tests:**

    ```sh
    make test
    ```

4. **Clean up build files:**

    ```sh
    make clean
    ```

### Using Docker

1. **Build Docker image:**

    ```sh
    make docker-build
    ```

2. **Run Docker container:**

    ```sh
    make docker-run
    ```

3. **Clean Docker container and image:**

    ```sh
    make docker-clean
    ```

## API Endpoints

### Add a Pokemon

- **URL:** `/pokemon`
- **Method:** `POST`
- **Request Body:**

    ```json
    {
        "id": 1,
        "name": "Bulbasaur",
        "type": "Grass/Poison",
        "height": 7,
        "weight": 69,
        "abilities": ["Overgrow", "Chlorophyll"]
    }
    ```

- **Response:** `201 Created` if successful

### Get Pokemon by ID

- **URL:** `/pokemon/id/{id}`
- **Method:** `GET`
- **Response:** `200 OK` with Pokemon data if found, `404 Not Found` if not found

### Get Pokemon by Name

- **URL:** `/pokemon/name/{name}`
- **Method:** `GET`
- **Response:** `200 OK` with Pokemon data if found, `404 Not Found` if not found

### Delete Pokemon by ID

- **URL:** `/pokemon/id/{id}`
- **Method:** `DELETE`
- **Response:** `200 OK` if successful, `404 Not Found` if not found

## Important Notes

- The default cache capacity is set to 100. You can adjust this in the `init` function in `handlers.go`.
- For development and testing, you can run the application locally using the Makefile. For deployment, you might prefer using Docker.

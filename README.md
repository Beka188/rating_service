# Politician Rating Service

This is a web service for managing ratings of politicians. It is built using the Gin framework in Go and uses SQLite as its database. The service provides endpoints to get all ratings, get a specific rating by ID, and update a rating.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Swagger Documentation](#swagger-documentation)
- [Contributing](#contributing)
- [License](#license)

## Features

- Retrieve all politician ratings.
- Retrieve a specific rating by ID.
- Update a politician's rating with actions such as upvote or downvote.
- Integrated with Swagger for API documentation.

## Installation

### Prerequisites

- Go 1.16+
- SQLite3

### Steps

1. Clone the repository:

    ```sh
    git clone https://github.com/Beka188/rating_service.git
    cd rating_service
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

[//]: # (3. Add your database file to the project directory and ensure it is named `foo.db`.)

3. Run the application:

    ```sh
    go run cmd/main.go
    ```

## Usage

### Running the Server

To run the server locally, execute the following command:

```sh
go run cmd/ main.go
```
The server will start on http://localhost:8080.

## API Endpoints
The following endpoints are available:


* GET /rating/
* GET /rating/{id}
* PUT /rating/{id}/{action}

### Swagger Documentation
This project includes Swagger documentation. You can access it at http://localhost:8080/swagger/index.html once the server is running.

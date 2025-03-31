# SWIFT codes API

This project is a Go-based application that parses SWIFT (BIC) data from a CSV file, stores it in a PostgreSQL database, and exposes the data through a RESTful API. The solution is containerized with Docker, allowing for a simple setup that includes both the production database and a separate test database.

## Table of Contents

- [SWIFT codes API](#swift-codes-api)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Project Structure](#project-structure)
  - [Prerequisites](#prerequisites)
  - [Installation \& Setup](#installation--setup)
    - [Local Setup](#local-setup)
    - [Docker Setup](#docker-setup)
  - [Usage (API Endpoints)](#usage-api-endpoints)
  - [Testing](#testing)
  - [Seed Data](#seed-data)

## Features

- **Data Parsing:** Reads a CSV file containing SWIFT data and processes each record.
- **Database Storage:** Stores parsed data in a PostgreSQL database with optimized schema and indexes.
- **RESTful API:** Exposes endpoints for:
  - Retrieving a single SWIFT code's details (with branches for headquarters).
  - Retrieving all SWIFT codes for a specific country.
  - Creating a new SWIFT code record.
  - Deleting a SWIFT code record.
- **Automated Seed:** On first run, if the production database is empty, the application will automatically seed it with data from the CSV file.
- **Test Environment:** Uses a separate test database for running unit and integration tests.
- **Containerization:** Fully containerized using Docker and Docker Compose for easy setup and deployment.

## Project Structure

```
swift-codes/
├── cmd/
│   ├── server/                  # Main server application (API)
│   │   └── main.go
│   └── import/                  # Import tool to seed the database (if used separately)
│       └── import.go
├── internal/
│   ├── db/                      # Database connection and CRUD operations
│   │   ├── db.go
│   │   └── db_test.go
│   ├── handlers/                # REST API endpoint implementations
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   ├── model/                   # Data model definitions
│   │   └── swift.go
│   └── parser/                  # CSV parsing logic
│       ├── parser.go
│       └── parser_test.go
├── data/                        
│   └── swiftcodes_data.csv      # CSV file for seeding data
├── entrypoint.sh                # Startup script for Docker that handles schema creation and seed import
├── Dockerfile                   # Dockerfile to build the application image
├── docker-compose.yml           # Docker Compose configuration for app, production DB, and test DB
├── go.mod                       # Go module file
└── README.md                    # This file
```

## Prerequisites

- [Go](https://golang.org/doc/install) 1.22
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Installation & Setup

### Local Setup
1. **Clone the Repository:**  
   ```git clone https://github.com/klark142/SWIFT-codes-API.git```  
   ```cd swift-codes```
2. **Set Up the Database Locally (Optional):**  
   Install PostgreSQL locally and create two databases: one for production (`swiftcodes`) and one for testing (`swiftcodes_test`).
3. **Set Environment Variables:**  
   Set the following variables in your terminal:
   - ```DB_CONN="host=localhost user=postgres password=secret dbname=swiftcodes sslmode=disable"```
   - `TEST_DB_CONN="host=localhost user=postgres password=secret dbname=swiftcodes_test sslmode=disable"`
4. **Run the Application:**  
   ```go run cmd/server/main.go```

### Docker Setup
1. **Clone the Repository:**  
   `git clone https://github.com/klark142/SWIFT-codes-API.git`  
   `cd swift-codes`
2. **Build and Run with Docker Compose:**  
   `docker-compose up --build`
   - This will build the application image and start three services:
     - **app:** The main application. It checks if the production database is empty and seeds it if necessary.
     - **db:** The production PostgreSQL database (`swiftcodes`).
     - **db_test:** The test PostgreSQL database (`swiftcodes_test`).

3. **Access the Application:**  
   The API will be available at [http://localhost:8080](http://localhost:8080).

## Usage (API Endpoints)
1. **GET /v1/swift-codes/{swiftCode}**  
   Retrieves details of a SWIFT code (if the record is a headquarters, branches are included).  
   Example: `curl http://localhost:8080/v1/swift-codes/AAISALTRXXX`

2. **GET /v1/swift-codes/country/{countryISO2code}**  
   Retrieves all SWIFT codes for a specific country.  
   Example: `curl http://localhost:8080/v1/swift-codes/country/BG`

3. **POST /v1/swift-codes**  
   Creates a new SWIFT code record.  
   Example:  
   ```curl -X POST http://localhost:8080/v1/swift-codes -H "Content-Type: application/json" -d '{"address": "Example Address", "bankName": "Example Bank", "countryISO2": "PL", "countryName": "POLAND", "isHeadquarter": true, "swiftCode": "EXAMPLEXXX"}'```

4. **DELETE /v1/swift-codes/{swift-code}**  
   Deletes a SWIFT code record.  
   Example: `curl -X DELETE http://localhost:8080/v1/swift-codes/EXAMPLEXXX`

## Testing
- The project includes unit and integration tests using a separate test database (`swiftcodes_test`) via the `TEST_DB_CONN` environment variable.
- **Running Tests Locally:**  
  Ensure your test database is set up and `TEST_DB_CONN` is configured, then run:  
  `go test ./...`
- **Running Tests with Docker:**  
  A separate service is provided in the Docker Compose configuration for tests. To run tests via Docker, execute:  
  `docker-compose run --rm app_test`
  This service builds an image with Go, runs `go test ./...`, and uses the test database configuration automatically.

## Seed Data
On the first run, the application checks if the production database is empty. If it is, it seeds the database by parsing a CSV file located in the `data` directory. This seeding process is performed only once.
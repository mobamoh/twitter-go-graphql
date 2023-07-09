# twitter-go-graphql

A simple Twitter clone project implemented in Go with GraphQL.

## Description

This project aims to create a basic Twitter-like application using Go programming language and GraphQL. It provides functionalities for creating tweets, following users, and retrieving tweet feeds.

## Libraries Used

- [PGX](https://github.com/jackc/pgx): PostgreSQL driver and toolkit for Go.
- [scany](https://github.com/georgysavva/scany): A lightweight ORM (Object-Relational Mapping) for Go.
- [gqlgen](https://github.com/99designs/gqlgen): A code generator for GraphQL servers in Go.
- [Dataloaden](https://github.com/vektah/dataloaden): A library for generating efficient data loaders in Go.
- [golang-migrate](https://github.com/golang-migrate/migrate): Database schema migration tool for Go applications.
- [Insomnia](https://insomnia.rest/): A powerful API testing and debugging tool.
- [Chi](https://github.com/go-chi/chi): Lightweight, idiomatic, and composable router for Go web applications.
- [PostgreSQL](https://www.postgresql.org/): A powerful open-source relational database management system.
- [godotenv](https://github.com/joho/godotenv): A Go package for reading environment variables from a `.env` file.
- [testify](https://github.com/stretchr/testify): A testing toolkit for Go.

## Setup and Installation

1. Clone the repository: `git clone https://github.com/your-username/twitter-go-graphql.git`
2. Install the required libraries: `go mod tidy`
3. Set up your PostgreSQL database and update the connection details in the `.env` file.
4. Run database migrations: `go run ./migrations/*.go up`
5. Start the application: `go run main.go`
6. Open Insomnia (or any other API testing tool) and import the provided Insomnia configuration file (`insomnia.json`) to test the GraphQL API.

## Usage

- Access the GraphQL API endpoint: `http://localhost:8080/graphql`
- Refer to the API documentation for available queries and mutations.
- Use Insomnia (or any other API testing tool) to send requests to the API.

## Testing

1. Run the tests: `go test ./...`

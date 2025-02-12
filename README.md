# Perpustakaan Golang

This project is a library management system built with Go, Gin, PostgreSQL, and Redis. It provides a RESTful API for managing books, members, and transactions, and includes Swagger documentation for easy integration.

## Features

- **Book Management:** Add, update, delete, and retrieve book information (title, author, ISBN, etc.).
- **Member Management:** (Implementation for member management is assumed based on the context, but code was not provided, so it is omitted from the instructions. If you provide member management code I will update the readme)
- **Transaction Management:** Borrow and return books, calculate fines.
- **Reporting:** Generate reports on borrowed and returned books.
- **User Authentication:** User registration, login, and token-based authentication using JWT (JSON Web Tokens)
- **Swagger Documentation:** Interactive API documentation is available via Swagger.
- **Dockerized Deployment:** Easy deployment using Docker Compose.
- **Redis Integration:** Uses Redis for JWT storage and potentially other caching or session management needs.
- **Image Uploads for Book and User:** Allows image uploading and retrieval using unique filenames and appropriate content types.

## API Documentation

You can access the interactive Swagger documentation after running the project by navigating to:

http://localhost:8080/swagger/index.html

## API Testing

You can import Postman collection and its environment to test this API app

## How to Run

### Using Docker Compose

`docker-compose build`

`docker compose up -d`

And then access the api from postman

The API base URL = http://127.0.0.1:9000/api/v1

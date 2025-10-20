# Chirpy API Documentation

Chirpy is a simple API for managing users and chirps (short messages similar to tweets). It provides endpoints for user registration, authentication, chirp creation, retrieval, and deletion.

## Base URL
The API runs on `http://localhost:8080` by default.

## Authentication
Some endpoints require authentication using JWT tokens. Include the token in the `Authorization` header as `Bearer <token>`.

Refresh tokens are used to obtain new access tokens.

## Endpoints

### Health Check
- **GET** `/api/healthz`
  - Description: Check if the server is running.
  - Response: `200 OK` with body "OK".

### User Management
- **POST** `/api/users`
  - Description: Create a new user.
  - Request Body: `{"email": "string", "password": "string"}`
  - Response: `201 Created` with User object.

- **PUT** `/api/users`
  - Description: Update user email and password. Requires authentication.
  - Request Body: `{"email": "string", "password": "string"}`
  - Response: `200 OK` with User object.

### Authentication
- **POST** `/api/login`
  - Description: Login and get access and refresh tokens.
  - Request Body: `{"email": "string", "password": "string"}`
  - Response: `200 OK` with User object including `token` and `refresh_token`.

- **POST** `/api/refresh`
  - Description: Get a new access token using refresh token.
  - Headers: `Authorization: Bearer <refresh_token>`
  - Response: `200 OK` with `{"token": "string"}`

- **POST** `/api/revoke`
  - Description: Revoke a refresh token.
  - Headers: `Authorization: Bearer <refresh_token>`
  - Response: `204 No Content`

### Chirps
- **POST** `/api/chirps`
  - Description: Create a new chirp. Requires authentication. Body is validated and cleaned.
  - Request Body: `{"body": "string"}`
  - Response: `201 Created` with Chirp object.

- **GET** `/api/chirps`
  - Description: Get all chirps. Optional query params: `author_id` (UUID), `sort` (asc/desc, default asc).
  - Response: `200 OK` with array of Chirp objects.

- **GET** `/api/chirps/{chirp_id}`
  - Description: Get a single chirp by ID.
  - Response: `200 OK` with Chirp object or `404 Not Found`.

- **DELETE** `/api/chirps/{chirp_id}`
  - Description: Delete a chirp. Requires authentication and ownership.
  - Response: `204 No Content` or error.

### Webhooks
- **POST** `/api/polka/webhooks`
  - Description: Upgrade user to Chirpy Red. Requires API key.
  - Request Body: `{"event": "user.upgraded", "data": {"user_id": "uuid"}}`
  - Response: `204 No Content`

### Admin
- **GET** `/admin/metrics`
  - Description: Get server metrics (visit count).
  - Response: `200 OK` with HTML page.

- **POST** `/admin/reset`
  - Description: Reset all users (dev only).
  - Response: `200 OK` with message.

## Data Models

### User
```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "string",
  "is_chirpy_red": boolean,
  "token": "string",  // optional
  "refresh_token": "string"  // optional
}
```

### Chirp
```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "body": "string",
  "user_id": "uuid"
}
```

## Setup
1. Set environment variables: `DB_URL`, `PLATFORM`, `SECRET_KEY`, `POLKA_KEY`.
2. Run `go run .` in the project directory.
3. The server starts on port 8080.

For more details, refer to the source code.

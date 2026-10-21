# Chirpy REST API Documentation

Welcome to the Chirpy REST API documentation. Chirpy is a microblogging platform service that provides user account management, posting and retrieval of short posts ("chirps"), subscription webhook integrations, and administrative monitoring.

## Environment Variables & Setup

Before running the Chirpy application server, you must set up your local environment configuration. The server uses `godotenv` to load environment variables from a `.env` file located in the root directory.

### `.env` File Setup

Create a `.env` file in the project root with the following variables:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=your-secret-jwt-signing-key
POLKA_KEY=your-polka-api-key
```

### Environment Configuration Parameters

| Variable | Required | Default / Mode | Description |
| :--- | :--- | :--- | :--- |
| `DB_URL` | **Yes** | — | PostgreSQL connection string used to connect to the primary database. |
| `PLATFORM` | **Yes** | `dev` | Set to `dev` to enable administrative reset endpoints (`POST /admin/reset`). In production, set to any value other than `dev` to lock down administrative reset capabilities. |
| `JWT_SECRET` | **Yes** | — | Secret key used to sign and validate Bearer JWT access tokens for authenticated user sessions. |
| `POLKA_KEY` | **Yes** | — | API key required to authenticate inbound webhooks from the Polka payment processor service (`POST /api/polka/webhooks`). |

---

## Overview & Architecture

### Base URL

By default, the HTTP server listens on port `8080`:

```
http://localhost:8080
```

### Request Headers

All endpoints accepting JSON payloads require:

```
Content-Type: application/json
```

Unless specified otherwise (e.g., text/html endpoints for admin metrics), response bodies are returned formatted as JSON.

### Authentication Schemes

The API supports two header-based authentication schemes depending on the endpoint:

1. **Bearer Token (JWT / Refresh Token)**
   * **Header:** `Authorization: Bearer <token>`
   * Used for authenticated user operations (creating/deleting chirps, updating profile) and session token management (refreshing access tokens, revoking refresh tokens).

2. **API Key Authentication**
   * **Header:** `Authorization: ApiKey <api_key>`
   * Used exclusively for third-party webhook verification (e.g., Polka payment notifications).

---

## Data Schemas

### User

Represents a user account in Chirpy.

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2026-10-21T10:00:00Z",
  "updated_at": "2026-10-21T10:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | string (UUID) | Unique identifier for the user. |
| `created_at` | string (ISO 8601) | Timestamp when the user account was created. |
| `updated_at` | string (ISO 8601) | Timestamp when the user account was last updated. |
| `email` | string | Email address associated with the user account. |
| `is_chirpy_red` | boolean | Subscription flag indicating if the user has upgraded to Chirpy Red. |

### Chirp

Represents a short message post (maximum 140 characters).

```json
{
  "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "created_at": "2026-10-21T10:05:00Z",
  "updated_at": "2026-10-21T10:05:00Z",
  "body": "This is a clean chirp post!",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | string (UUID) | Unique identifier for the chirp. |
| `created_at` | string (ISO 8601) | Timestamp when the chirp was posted. |
| `updated_at` | string (ISO 8601) | Timestamp when the chirp was last updated. |
| `body` | string | Message body. Profane words (`kerfuffle`, `sharbert`, `fornax`) are automatically sanitized to `****`. |
| `user_id` | string (UUID) | Identifier of the authoring user. |

---

## Endpoint Reference

### Health & Admin

#### 1. Health Check
Checks if the server is healthy and operational.

* **HTTP Method:** `GET`
* **Path:** `/api/healthz`
* **Authentication:** None

##### Responses
* **`200 OK`** - Server is healthy.
  * **Content-Type:** `text/plain; charset=utf-8`
  * **Body:** `Endpoint is OK`

#### 2. Get Server Metrics
Returns an HTML page showing total file server hit count.

* **HTTP Method:** `GET`
* **Path:** `/admin/metrics`
* **Authentication:** None

##### Responses
* **`200 OK`**
  * **Content-Type:** `text/html; charset=utf-8`
  * **Body:**
    ```html
    <html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited 42 times!</p></body></html>
    ```

#### 3. Reset Application State
Resets the file server hit counter to `0` and deletes all user accounts from the database. **Only allowed when `PLATFORM=dev`**.

* **HTTP Method:** `POST`
* **Path:** `/admin/reset`
* **Authentication:** None

##### Responses
* **`200 OK`** - Reset successfully executed.
  * **Content-Type:** `text/plain; charset=utf-8`
  * **Body:** `All users deleted and hits reset to 0`
* **`403 Forbidden`** - Attempted invocation while `PLATFORM` is not set to `dev`.
  * **Body:** `Unauthorized Access attempt, FBI notified!`
* **`500 Internal Server Error`** - Failed to reset database records.

---

### Authentication & Users

#### 1. Register User
Creates a new user account with an email and password.

* **HTTP Method:** `POST`
* **Path:** `/api/users`
* **Authentication:** None

##### Request Body
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

##### Responses
* **`201 Created`** - User account successfully created.
  ```json
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2026-10-21T10:00:00Z",
    "updated_at": "2026-10-21T10:00:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false
  }
  ```
* **`400 Bad Request`** - Invalid or malformed JSON request body.
* **`500 Internal Server Error`** - Failed to hash password or persist user in database.

#### 2. User Login
Authenticates user credentials and issues a JWT access token and a refresh token (valid for 60 days).

* **HTTP Method:** `POST`
* **Path:** `/api/login`
* **Authentication:** None

##### Request Body
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

##### Responses
* **`200 OK`** - Successfully authenticated.
  ```json
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2026-10-21T10:00:00Z",
    "updated_at": "2026-10-21T10:00:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "608d0e70a2..."
  }
  ```
* **`400 Bad Request`** - Invalid JSON request body.
* **`401 Unauthorized`** - Incorrect email or password.
* **`500 Internal Server Error`** - Failed to issue access or refresh token.

#### 3. Update User Credentials
Updates the email address and password for the currently authenticated user.

* **HTTP Method:** `PUT`
* **Path:** `/api/users`
* **Authentication:** `Bearer <JWT_TOKEN>`

##### Request Body
```json
{
  "email": "newemail@example.com",
  "password": "newsecurepassword123"
}
```

##### Responses
* **`200 OK`** - User credentials updated successfully.
  ```json
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2026-10-21T10:00:00Z",
    "updated_at": "2026-10-21T10:30:00Z",
    "email": "newemail@example.com",
    "is_chirpy_red": false
  }
  ```
* **`400 Bad Request`** - Invalid JSON request payload.
* **`401 Unauthorized`** - Missing, invalid, or expired JWT token.
* **`500 Internal Server Error`** - Failed to hash new password or update database record.

#### 4. Refresh Access Token
Generates a new JWT access token using a valid Bearer refresh token.

* **HTTP Method:** `POST`
* **Path:** `/api/refresh`
* **Authentication:** `Bearer <REFRESH_TOKEN>`

##### Responses
* **`200 OK`** - Access token refreshed.
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
* **`401 Unauthorized`** - Invalid, missing, expired, or revoked refresh token.
* **`500 Internal Server Error`** - Failed to issue new JWT token.

#### 5. Revoke Refresh Token
Revokes a refresh token so it cannot be used to generate new access tokens.

* **HTTP Method:** `POST`
* **Path:** `/api/revoke`
* **Authentication:** `Bearer <REFRESH_TOKEN>`

##### Responses
* **`204 No Content`** - Refresh token revoked successfully.
* **`401 Unauthorized`** - Missing or invalid refresh token.

---

### Chirps

#### 1. Create Chirp
Creates a new chirp post for the authenticated user. Posts are capped at 140 characters. Words `kerfuffle`, `sharbert`, and `fornax` are replaced with `****`.

* **HTTP Method:** `POST`
* **Path:** `/api/chirps`
* **Authentication:** `Bearer <JWT_TOKEN>`

##### Request Body
```json
{
  "body": "Hello world! This is a kerfuffle test."
}
```

##### Responses
* **`201 Created`** - Chirp created successfully.
  ```json
  {
    "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "created_at": "2026-10-21T10:05:00Z",
    "updated_at": "2026-10-21T10:05:00Z",
    "body": "Hello world! This is a **** test.",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
  ```
* **`400 Bad Request`** - Body exceeds 140 characters or invalid JSON format.
* **`401 Unauthorized`** - Missing or invalid JWT access token.
* **`500 Internal Server Error`** - Failed to persist chirp to the database.

#### 2. Get All Chirps
Retrieves a list of all chirps, optional filtered by author and sorted by creation date.

* **HTTP Method:** `GET`
* **Path:** `/api/chirps`
* **Authentication:** None

##### Query Parameters
| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `author_id` | string (UUID) | No | None | Filters chirps authored by a specific user ID. |
| `sort` | string | No | `asc` | Controls sort ordering by `created_at`. Accepted values: `asc`, `desc`. |

##### Responses
* **`200 OK`**
  ```json
  [
    {
      "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
      "created_at": "2026-10-21T10:05:00Z",
      "updated_at": "2026-10-21T10:05:00Z",
      "body": "Hello world!",
      "user_id": "123e4567-e89b-12d3-a456-426614174000"
    }
  ]
  ```
* **`400 Bad Request`** - `author_id` parameter is not a valid UUID format.
* **`500 Internal Server Error`** - Database retrieval error.

#### 3. Get Chirp by ID
Retrieves a single chirp by its unique UUID.

* **HTTP Method:** `GET`
* **Path:** `/api/chirps/{chirpID}`
* **Authentication:** None

##### URL Parameters
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `chirpID` | string (UUID) | Yes | Unique identifier of the chirp. |

##### Responses
* **`200 OK`**
  ```json
  {
    "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "created_at": "2026-10-21T10:05:00Z",
    "updated_at": "2026-10-21T10:05:00Z",
    "body": "Hello world!",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
  ```
* **`400 Bad Request`** - `chirpID` path parameter is not a valid UUID.
* **`404 Not Found`** - No chirp exists with the specified ID.

#### 4. Delete Chirp
Deletes a chirp by ID. **Only the original author of the chirp can delete it.**

* **HTTP Method:** `DELETE`
* **Path:** `/api/chirps/{chirpID}`
* **Authentication:** `Bearer <JWT_TOKEN>`

##### URL Parameters
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `chirpID` | string (UUID) | Yes | Unique identifier of the chirp to delete. |

##### Responses
* **`204 No Content`** - Chirp deleted successfully.
* **`400 Bad Request`** - Invalid UUID path parameter.
* **`401 Unauthorized`** - Missing or invalid JWT access token.
* **`403 Forbidden`** - Authenticated user is not the author of the chirp.
* **`404 Not Found`** - Chirp not found.
* **`500 Internal Server Error`** - Error executing deletion query.

---

### Webhooks

#### 1. Polka Webhook Integration
Receives external webhook events from Polka. When a `user.upgraded` event is received, the target user account's `is_chirpy_red` status is updated to `true`.

* **HTTP Method:** `POST`
* **Path:** `/api/polka/webhooks`
* **Authentication:** `ApiKey <POLKA_KEY>`

##### Request Body
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

##### Responses
* **`204 No Content`** - Webhook processed successfully or ignored (if event type is not `user.upgraded`).
* **`400 Bad Request`** - Invalid request body structure.
* **`401 Unauthorized`** - Missing or invalid Polka API Key.
* **`404 Not Found`** - User ID specified in webhook payload does not exist in the database.
* **`500 Internal Server Error`** - Error updating user record.
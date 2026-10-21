# Chirpy REST API Documentation

Welcome to the Chirpy REST API documentation. Chirpy is a microblogging platform service that provides user account management, posting and retrieval of short posts ("chirps"), subscription webhook integrations, and administrative monitoring.

---

## Overview & Architecture

### Base URL
By default, the server runs locally at:
```text
http://localhost:8080
```

### Content Type
All request bodies sent to JSON endpoints must include the header:
```http
Content-Type: application/json
```
Unless specified otherwise (e.g., text/html endpoints for admin metrics), response bodies are returned formatted as JSON.

### Authentication Schemes
The API supports two header-based authentication schemes depending on the endpoint:

1. **Bearer Token (JWT / Refresh Token)**
   - Header format: `Authorization: Bearer <token>`
   - Used for authenticated user endpoints (creating/deleting chirps, updating profile) and session token management (refresh, revoke).
2. **API Key Authentication**
   - Header format: `Authorization: ApiKey <api_key>`
   - Used exclusively for third-party webhook verification (e.g., Polka integration).

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
| `created_at` | string (ISO 8601) | Timestamp when user was created. |
| `updated_at` | string (ISO 8601) | Timestamp when user was last updated. |
| `email` | string | Email address of the user. |
| `is_chirpy_red` | boolean | Status indicating if the user has upgraded to Chirpy Red. |

---

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
| `created_at` | string (ISO 8601) | Timestamp when chirp was created. |
| `updated_at` | string (ISO 8601) | Timestamp when chirp was last updated. |
| `body` | string | Post content. Profanity (`kerfuffle`, `sharbert`, `fornax`) is automatically replaced with `****`. |
| `user_id` | string (UUID) | ID of the user who authored the chirp. |

---

## Endpoint Reference

---

### Health & Admin

#### 1. Health Check
Checks if the server is healthy and running.

* **HTTP Method:** `GET`
* **Path:** `/api/healthz`
* **Authentication:** None

##### Responses
* **`200 OK`** - Server is healthy.
  * **Content-Type:** `text/plain; charset=utf-8`
  * **Response Body:**
    ```text
    Endpoint is OK
    ```

---

#### 2. Get Server Metrics
Returns an HTML page detailing server hit count.

* **HTTP Method:** `GET`
* **Path:** `/admin/metrics`
* **Authentication:** None

##### Responses
* **`200 OK`**
  * **Content-Type:** `text/html; charset=utf-8`
  * **Response Body:**
    ```html
    <html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited 42 times!</p></body></html>
    ```

---

#### 3. Reset Application State
Resets the server file hit counter and deletes all user accounts from the database. *Only allowed when `PLATFORM=dev`.*

* **HTTP Method:** `POST`
* **Path:** `/admin/reset`
* **Authentication:** None

##### Responses
* **`200 OK`** - Reset successful.
  * **Content-Type:** `text/plain; charset=utf-8`
  * **Response Body:**
    ```text
    All users deleted and hits reset to 0
    ```
* **`403 Forbidden`** - Platform is not configured to `dev`.
  * **Response Body:**
    ```text
    Unauthorized Access attempt, FBI notified!
    ```
* **`500 Internal Server Error`** - Failed to purge database records.

---

### Authentication & Users

#### 1. Register User
Creates a new user account.

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
* **`201 Created`** - User successfully created.
  ```json
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2026-10-21T10:00:00Z",
    "updated_at": "2026-10-21T10:00:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false
  }
  ```
* **`400 Bad Request`** - Invalid JSON request body.
* **`500 Internal Server Error`** - Password hashing or user creation failed.

---

#### 2. User Login
Authenticates a user and returns access JWT and refresh tokens.

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
* **`200 OK`** - Authentication successful.
  ```json
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2026-10-21T10:00:00Z",
    "updated_at": "2026-10-21T10:00:00Z",
    "email": "user@example.com",
    "is_chirpy_red": false,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "608d0e70a2...60daysvalid"
  }
  ```
* **`400 Bad Request`** - Invalid JSON request payload.
* **`401 Unauthorized`** - Invalid email or password.
* **`500 Internal Server Error`** - Error issuing authentication tokens.

---

#### 3. Update User Credentials
Updates the email and password for the authenticated user.

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
* **`200 OK`** - Account details updated successfully.
  ```json
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2026-10-21T10:00:00Z",
    "updated_at": "2026-10-21T10:30:00Z",
    "email": "newemail@example.com",
    "is_chirpy_red": false
  }
  ```
* **`400 Bad Request`** - Invalid request body.
* **`401 Unauthorized`** - Missing or invalid JWT access token.
* **`500 Internal Server Error`** - Failed to process update.

---

#### 4. Refresh Access Token
Issues a new JWT access token using a valid, non-revoked refresh token.

* **HTTP Method:** `POST`
* **Path:** `/api/refresh`
* **Authentication:** `Bearer <REFRESH_TOKEN>`

##### Responses
* **`200 OK`** - Token successfully refreshed.
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
* **`401 Unauthorized`** - Missing, invalid, expired, or revoked refresh token.
* **`500 Internal Server Error`** - Token generation failed.

---

#### 5. Revoke Refresh Token
Revokes a refresh token, rendering it unusable for future access token requests.

* **HTTP Method:** `POST`
* **Path:** `/api/revoke`
* **Authentication:** `Bearer <REFRESH_TOKEN>`

##### Responses
* **`204 No Content`** - Refresh token revoked successfully.
* **`401 Unauthorized`** - Token missing or invalid.

---

### Chirps

#### 1. Create Chirp
Posts a new chirp on behalf of the authenticated user. Post body max length is 140 characters. Automatically censors profanity (`kerfuffle`, `sharbert`, `fornax`).

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
* **`201 Created`** - Chirp created.
  ```json
  {
    "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "created_at": "2026-10-21T10:05:00Z",
    "updated_at": "2026-10-21T10:05:00Z",
    "body": "Hello world! This is a **** test.",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
  ```
* **`400 Bad Request`** - Invalid request body or chirp exceeds 140 characters.
* **`401 Unauthorized`** - Missing or invalid JWT authentication.
* **`500 Internal Server Error`** - Database insertion failed.

---

#### 2. Get All Chirps
Retrieves a list of chirps, optionally filtered by author ID and sorted by creation time.

* **HTTP Method:** `GET`
* **Path:** `/api/chirps`
* **Authentication:** None

##### Query Parameters
| Parameter | Type | Required | Default | Description |
| :--- | :--- | :--- | :--- | :--- |
| `author_id` | string (UUID) | No | None | Filters chirps written by a specific user UUID. |
| `sort` | string | No | `asc` | Sort direction based on `created_at`. Values: `asc`, `desc`. |

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
* **`400 Bad Request`** - Provided `author_id` is not a valid UUID format.
* **`500 Internal Server Error`** - Database query failure.

---

#### 3. Get Chirp by ID
Retrieves a specific chirp by its unique UUID.

* **HTTP Method:** `GET`
* **Path:** `/api/chirps/{chirpID}`
* **Authentication:** None

##### URL Parameters
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `chirpID` | string (UUID) | Yes | Unique UUID identifier of the chirp. |

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
* **`400 Bad Request`** - Invalid UUID path parameter.
* **`404 Not Found`** - No chirp exists with the provided ID.

---

#### 4. Delete Chirp
Deletes a chirp by ID. Only the original author of the chirp can delete it.

* **HTTP Method:** `DELETE`
* **Path:** `/api/chirps/{chirpID}`
* **Authentication:** `Bearer <JWT_TOKEN>`

##### URL Parameters
| Parameter | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| `chirpID` | string (UUID) | Yes | Unique UUID identifier of the chirp to delete. |

##### Responses
* **`204 No Content`** - Chirp deleted successfully.
* **`400 Bad Request`** - Invalid UUID path parameter.
* **`401 Unauthorized`** - Missing or invalid JWT authentication token.
* **`403 Forbidden`** - Authenticated user is not the author of the chirp.
* **`404 Not Found`** - Chirp not found.
* **`500 Internal Server Error`** - Deletion operation failed.

---

### Webhooks

#### 1. Polka Webhook Event
Processes events sent by the third-party Polka service. Upgrades a user account to **Chirpy Red** when receiving `user.upgraded` events.

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
* **`204 No Content`** - Processed successfully (or non-handled event ignored).
* **`400 Bad Request`** - Invalid request body.
* **`401 Unauthorized`** - Missing or invalid Polka API Key.
* **`404 Not Found`** - User specified in webhook payload does not exist.
* **`500 Internal Server Error`** - Error executing user account upgrade.
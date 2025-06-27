# Chirpy

Chirpy is a twitter-style social network with user registration.

I made this project to learn about HTTP servers and databases.


## Features

- User registration and login
- Password hashing with bcrypt
- JWT authentication and refresh token system
- Create, list, and delete chirps
- Polkadot style webhook
- Static file serving

## Getting Started

### Prerequisites

- Go 1.24.3+
- PostgreSQL
- Goose (for database migrations)
- sqlc (for generating type-safe DB code)

### Setup

1. Clone the repo

```
git clone https://github.com/npayetteraynauld/Chirpy
```

2. Create a .env file

```
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=chirpy
SECRET=your_jwt_secret
POLKA_KEY=your_polka_webhook_key
```

3. Setup database with provided SQL migrations using goose

```
goose postgres <DB_URL> up
```

4. Run the server

```
go run .
```

server will run on http://localhost:8080

## Api

### Health

```
GET /api/healthz
```
Returns 200 OK and a basic healthy response.

### Users

#### Register

```
POST /api/users
```

##### Request body:

```
{
  "email": "",
  "password": ""
}
```

##### Response:

```
{
  "id": 123,
  "created_at": "",
  "updated_at": "",
  "email": "",
  "is_chirpy_red": bool
}
```

#### Update Profile

```
PUT /api/users
```

Requires authorization: Bearer <access_token>

##### Request body:

```
{
  "email": "",
  "password": ""
}
```

##### Response:

```
{
  "id": 123,
  "created_at": "",
  "updated_at": "",
  "email": "",
  "is_chirpy_red": bool
}
```

### Auth

#### Login

```
POST /api/login
```

Returns access and refresh token

##### Request body:

```
{
  "email": "",
  "password": ""
}
```

##### Response: 

```
{
  "id": 123,
  "created_at": "",
  "updated_at": "",
  "email": "",
  "token": ""
  "refresh_token": ""
  "is_chirpy_red": bool
}
```

#### Refresh

```
POST /api/refresh
```

Requires authorization: Bearer <access_token>

##### Response

```
{
  "token": asd
}
```


#### Revoke

```
POST /api/revoke
```

Revokes refresh token
Requires authorization: Bearer <access_token>

### Chirps

#### List Chirps

```
GET /api/chirps?author_id=<uuid>?sort=<asc/desc>
```

Returns chirps, optionally filtered by author_id, and optionally sorted by asc or desc created_at timestamp

#### Create Chirp

```
POST /api/chirps
```

Requires JWT

##### Request body:

```
{
  "body": "This is a new chirp!"
}
```

#### Chirp by ID

```
GET /api/chirps/{chirpID}
```

```
DELETE /api/chirps/{Chirpy}
```

Deleting requires JWT

### Polkadot Webhooks

```
POST /api/polka/webhooks
```

Expects header: ApiKey <webhook_key>

### Admin

#### Metrics

```
GET /admin/metrics
```

Shows file-serving hit count

#### Reset

```
POST /admin/reset
```

Resets database

# Learn Go Event API

Backend API Go untuk event ticketing: auth, venue, guest/speaker, event, order, ticket, payment Midtrans, Redis event bus, dan PostgreSQL persistence.

## Stack

- Go 1.24
- Gin HTTP router
- Cobra CLI
- GORM + PostgreSQL
- Redis
- Midtrans Core API
- Viper config

## Run dengan Docker

```bash
docker compose up --build
```

API berjalan di:

```text
http://localhost:8080/api/v1
```

## Run lokal

Siapkan PostgreSQL dan Redis, lalu set environment variable atau `.env`:

```env
APP_ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=learngo
JWT_SECRET_KEY=replace-with-long-random-secret
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
MIDTRANS_SERVER_KEY=
MIDTRANS_CLIENT_KEY=
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=
```

Jalankan server:

```bash
go run . serve
```

Jalankan migration manual:

```bash
go run . migrate
```

Jalankan compile/test check:

```bash
go test ./...
```

## Auth

Login mengembalikan token dan juga memasang cookie `jwt_token` dengan:

- `HttpOnly`
- `SameSite=Lax`
- `Secure=true` saat `APP_ENV=production`
- expiry 24 jam

JWT memakai HS256 dan divalidasi dengan issuer/audience. Protected route memerlukan user yang sudah verified.

## Request ID dan logging

Setiap response memiliki header:

```text
X-Request-ID
```

Jika client mengirim `X-Request-ID`, nilai itu dipropagasi. Jika tidak, server membuat request ID baru. Request log JSON menyertakan method, path, route, status, duration, client IP, user agent, request ID, dan user ID jika tersedia.

## Rate limit

Endpoint sensitif diberi rate limit berbasis Redis dan key per user/IP.

Contoh endpoint yang dilimit:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/verify-otp`
- `POST /api/v1/auth/login`
- `POST /api/v1/orders/`
- `POST /api/v1/payments/`
- `PATCH /api/v1/payments/:id/status`
- `POST /api/v1/payments/midtrans-notification`

Response limit memakai status `429` dan header:

```text
X-RateLimit-Limit
X-RateLimit-Remaining
Retry-After
```

## Order dan payment lifecycle

Order status:

```text
PENDING -> PAID
PENDING -> CANCELLED
```

Payment status:

```text
PENDING -> SUCCESS
PENDING -> FAILED
SUCCESS -> REFUNDED
```

Transisi invalid, seperti `SUCCESS -> PENDING`, ditolak/diabaikan.

Payment notification Midtrans:

- signature diverifikasi
- payload wajib memiliki `order_id`, `transaction_status`, dan `transaction_id`
- notifikasi duplikat bersifat idempotent
- update payment dan order dilakukan dalam database transaction

## API docs

OpenAPI draft tersedia di:

```text
docs/openapi.yaml
```

Postman collection juga tersedia di:

```text
postman_collection.json
```

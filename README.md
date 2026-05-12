# classifier-ticket-department-backend

Predict Ticket Department Backend API (Gin + PostgreSQL).

This API supports:

1. Company authentication (register/login)
2. Company form and ticket workflows
3. Department management
4. Internal HMAC-secured endpoints


## Authentication

This project has 2 authentication styles:

1. Bearer JWT
	 - Header: `Authorization: Bearer <token>`
	 - Token source: `POST /api/v1/login`
2. HMAC
	 - Header: `X-HMAC-Signature: sha256=<signature>`
	 - Used for internal protected routes

## API Documentation

Base URL: `http://localhost:8888`

### Company

#### POST /api/v1/register

Register company account.

Request body:

```json
{
	"email": "admin@company.com",
	"password": "123456"
}
```

Status: `201`, `400`, `409`, `500`

#### POST /api/v1/login

Login and get JWT token.

Request body:

```json
{
	"email": "admin@company.com",
	"password": "123456"
}
```

Status: `200`, `400`, `401`

### Company Form (JWT)

#### POST /api/v1/links

Create company form link metadata.

Headers:

```text
Authorization: Bearer <token>
```

Request body:

```json
{
	"company_id": "uuid-company-id"
}
```

Status: `201`, `400`, `500`

#### GET /api/v1/links

Get link form by company id.

Headers:

```text
Authorization: Bearer <token>
```

Status: `200`, `400`, `500`

### Departments

#### POST /api/v1/departments (JWT)

Add departments to a company.

Headers:

```text
Authorization: Bearer <token>
```

Request body:

```json
{
	"department_name": ["IT", "HR", "Support"]
}
```

Status: `201`, `400`, `500`

#### GET /api/v1/internal/departments/company/{company_id} (HMAC)

Get departments by company id.

Headers:

```text
X-HMAC-Signature: sha256=<signature>
```

Status: `200`, `400`, `500`

### Forms

#### POST /forms/submit

Submit a form.

Request body:

```json
{
	"title": "Cannot access VPN",
	"description": "VPN disconnected every 5 minutes",
	"link_id": "uuid-link-id"
}
```

Status: `201`, `400`, `500`

#### GET /api/v1/forms

Get all forms by company id.

Status: `200`, `400`, `500`

#### GET /api/v1/forms/per-day?date=YYYY-MM-DD

Get forms by company id per day.

Query:

1. `date` optional, format `YYYY-MM-DD`
2. default is yesterday when omitted

Status: `200`, `400`, `500`

### Tickets

#### GET /api/v1/tickets

Get all tickets by company id.

Status: `200`, `400`, `500`

#### POST /api/v1/tickets (JWT)

Create single ticket.

Headers:

```text
Authorization: Bearer <token>
```

Request body:

```json
{
	"form_id": "uuid-form-id",
	"department_id": "uuid-department-id",
	"title": "VPN issue",
	"description": "User cannot connect to VPN",
	"message": "Need urgent support",
	"priority": "high",
	"status": "success"
}
```

Status: `200`, `400`, `500`

#### POST /api/v1/tickets/bulk (HMAC)

Create multiple tickets in one request.

Headers:

```text
X-HMAC-Signature: sha256=<signature>
```

Request body:

```json
[
	{
		"form_id": "uuid-form-id",
		"department_id": "uuid-department-id",
		"title": "VPN issue",
		"description": "User cannot connect to VPN",
		"message": "Need urgent support",
		"priority": "high",
		"status": "success"
	}
]
```

Status: `201`, `400`, `500`



## Troubleshooting

1. `docker-compose up` fails: make sure Docker Desktop is running.
2. Database connection error: verify `DB_PASS`, `DB_HOST`, and `DB_PORT` in `.env`.
3. Unauthorized on JWT routes: confirm `Authorization: Bearer <token>` is set.
4. HMAC route rejected: verify `X-HMAC-Signature` format and shared secret.

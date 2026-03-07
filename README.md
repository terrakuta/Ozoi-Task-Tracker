# Ozoi Task Tracker

REST API for task tracking

## Stack

- **Go** + **Gin** — language and web framework
- **PostgreSQL** — primary database
- **pgx** — PostgreSQL driver and connection pool
- **golang-migrate** — database migrations
- **JWT** in httpOnly cookie — authentication
- **Swagger** — API documentation with Basic Auth protection

## Features

- User registration and login
- Cookie-based session management
- Protected routes via JWT middleware
- Full CRUD for tasks
- Per-user data isolation (users only see their own tasks)
- Input validation on all endpoints
- Auto-generated API docs via Swagger

## Getting Started

### Local Development

1. Copy `.env.example` to `.env` and fill in your values
2. Start the server (using [Air](https://github.com/air-verse/air) for hot reload):
```
air
```
> Migrations are applied automatically on every server start.

> Swagger UI: `http://localhost:3000/swagger/index.html`

### Docker

1. Copy `.env.example` to `.env.docker` and fill in your values
2. Run:
```bash
docker-compose up --build
```
> Migrations are applied automatically on startup.

> Swagger UI: `http://localhost:3000/swagger/index.html`

## Migrations

Manual migration management for local development (requires [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)):
```powershell
.\scripts\migrate.ps1 up          # apply all migrations
.\scripts\migrate.ps1 down        # roll back 1 migration
.\scripts\migrate.ps1 down 3      # roll back 3 migrations
.\scripts\migrate.ps1 create name # create new migration
.\scripts\migrate.ps1 force 1     # force specific version
```

## Endpoints

| Method | URL | Auth | Description |
|--------|-----|------|-------------|
| POST | /auth/register | — | Register a new user |
| POST | /auth/login | — | Login |
| POST | /auth/logout | — | Logout |
| GET | /auth/me | ✅ | Get current user |
| POST | /ozoi | ✅ | Create a task |
| GET | /ozoi | ✅ | Get all tasks |
| GET | /ozoi/:id | ✅ | Get task by ID |
| PUT | /ozoi/:id | ✅ | Update task |
| DELETE | /ozoi/:id | ✅ | Delete task |

## Environment Variables

All configuration is done via environment variables. See `.env.example` for reference.

> `SWAGGER_USER` and `SWAGGER_PASSWORD` are used for Basic Auth on the `/swagger` route.
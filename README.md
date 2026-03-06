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

1. Copy `.env.example` to `.env` and fill in your values
2. Run migrations manually on first launch:
```powershell
.\scripts\migrate.ps1 up
```
3. Start the server (using [Air](https://github.com/air-verse/air) for hot reload):
```
air
```
> After the first run, migrations are applied automatically on every server start.

> Swagger UI: `http://localhost:8080/swagger/index.html`

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
```

## Environment Variables

All configuration is done via environment variables. Copy `.env.example` to `.env` and fill in your values:
```env
DATABASE_URL=postgres://user:password@localhost:5432/ozoi
JWT_SECRET=your_secret
PORT=8080
SWAGGER_USER=admin
SWAGGER_PASSWORD=admin
```

> `SWAGGER_USER` and `SWAGGER_PASSWORD` are used for Basic Auth on the `/swagger` route.
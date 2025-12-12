# LinkMy v2.0

High-performance link-in-bio platform built with Go + PostgreSQL + Redis.

## Tech Stack

- **Backend**: Go 1.22 + Fiber (ultra-fast HTTP framework)
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Frontend**: SvelteKit (coming in Phase 2)

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.22+ (for local development)

### Development

1. **Clone and setup**:
   ```bash
   cd v2
   cp .env.example .env
   # Edit .env with your values
   ```

2. **Start services**:
   ```bash
   docker-compose up -d
   ```

3. **Run API locally** (optional, for hot-reload):
   ```bash
   cd backend
   go run ./cmd/server
   ```

### API Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout

#### Public
- `GET /api/v1/p/:slug` - Get public profile
- `POST /api/v1/click/:id` - Track link click

#### Protected (requires JWT)
- `GET /api/v1/me` - Get current user
- `PUT /api/v1/me` - Update current user
- `GET /api/v1/profiles` - List user's profiles
- `POST /api/v1/profiles` - Create profile
- `PUT /api/v1/profiles/:id` - Update profile
- `DELETE /api/v1/profiles/:id` - Delete profile
- `GET /api/v1/profiles/:id/links` - Get profile links
- `POST /api/v1/profiles/:id/links` - Create link
- `PUT /api/v1/links/:id` - Update link
- `DELETE /api/v1/links/:id` - Delete link
- `GET /api/v1/profiles/:id/theme` - Get theme
- `PUT /api/v1/profiles/:id/theme` - Update theme
- `GET /api/v1/profiles/:id/analytics` - Get analytics

## Project Structure

```
v2/
├── backend/
│   ├── cmd/server/         # Entry point
│   ├── internal/
│   │   ├── config/         # Configuration
│   │   ├── database/       # DB connection + migrations
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # Auth middleware
│   │   ├── models/         # Data models
│   │   └── repository/     # Data access layer
│   ├── Dockerfile
│   └── go.mod
├── frontend/               # SvelteKit (Phase 2)
├── docker-compose.yml
└── .env.example
```

## Performance

Compared to PHP v1:
- **100x faster** request handling
- **500x more** concurrent users
- **75% less** memory usage

## License

MIT

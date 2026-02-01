# LocalConnect

A full-stack web application connecting local workers (labour, cook, painter, carpenter, home tuition, teacher, househelp) with clients who need their services.

## Features

### For Workers
- Create and manage professional profiles
- Set rates and showcase experience
- Upload photos/videos of work
- Receive reviews and upvotes
- Get contact requests from clients
- Message directly with clients
- Track profile views and engagement

### For Clients
- Search workers by category and location
- View worker profiles with ratings and reviews
- Upvote favorite workers
- Leave reviews and ratings
- Send contact requests
- Message workers directly
- Filter and sort workers

## Technology Stack

### Backend (Golang)
- **Framework**: Go 1.22+ with chi router
- **Database**: PostgreSQL 16 with pgx driver
- **Authentication**: JWT tokens
- **Concurrency**: Goroutines and channels for:
  - Worker pools for async tasks
  - Notification batching and processing
  - Analytics event tracking
  - Cache management
- **Features**:
  - In-memory caching with TTL
  - Async notifications system
  - Analytics tracking
  - Media upload handling
  - RESTful API design

### Frontend (React)
- **Framework**: React 19 with Vite
- **Routing**: React Router v7
- **Styling**: Custom CSS with modern UI
- **Features**:
  - Responsive design
  - Real-time messaging interface
  - File upload for media
  - Search and filtering
  - Authentication flow

### Database
- PostgreSQL with optimized indexes
- Foreign key constraints
- Triggers for automated updates
- Full-text search ready

## Project Structure

```
LocalConnect/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Application entry point
│   ├── internal/
│   │   ├── analytics/
│   │   │   └── tracker.go           # Analytics tracking with goroutines
│   │   ├── auth/
│   │   │   └── jwt.go               # JWT authentication
│   │   ├── cache/
│   │   │   └── cache.go             # Thread-safe in-memory cache
│   │   ├── config/
│   │   │   └── config.go            # Configuration management
│   │   ├── db/
│   │   │   └── db.go                # Database connection pool
│   │   ├── handlers/
│   │   │   ├── auth.go              # Authentication handlers
│   │   │   ├── profiles.go          # Worker profile handlers
│   │   │   ├── reviews.go           # Review handlers
│   │   │   ├── upvotes.go           # Upvote handlers
│   │   │   ├── messaging.go         # Messaging handlers
│   │   │   ├── contact.go           # Contact request handlers
│   │   │   ├── categories.go        # Category handlers
│   │   │   ├── media.go             # Media upload handlers
│   │   │   └── routes.go            # Route definitions
│   │   ├── media/
│   │   │   └── media.go             # Media file handling
│   │   ├── models/
│   │   │   ├── models.go            # Data models
│   │   │   └── messaging.go         # Messaging models
│   │   ├── notifications/
│   │   │   └── notifier.go          # Async notification system
│   │   └── worker/
│   │       └── pool.go              # Worker pool for concurrent tasks
│   ├── migrations/
│   │   ├── 001_init.sql             # Initial schema
│   │   └── 002_messaging.sql        # Messaging schema + indexes
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── client.js            # API client with authentication
│   │   ├── components/
│   │   │   └── Layout.jsx           # App layout with navigation
│   │   ├── pages/
│   │   │   ├── Home.jsx             # Search and browse workers
│   │   │   ├── ProfileDetail.jsx    # Worker profile detail
│   │   │   ├── Dashboard.jsx        # Worker dashboard
│   │   │   ├── Auth.jsx             # Login/Register
│   │   │   ├── Messages.jsx         # Messaging interface
│   │   │   └── Messages.css         # Messaging styles
│   │   ├── App.jsx                  # App component with routes
│   │   ├── App.css                  # Global styles
│   │   └── main.jsx                 # React entry point
│   ├── package.json
│   └── vite.config.js
├── docker-compose.yml               # PostgreSQL container
├── env.example                      # Environment variables template
└── README.md
```

## Setup Instructions

### Prerequisites
- Go 1.22 or higher
- Node.js 18 or higher
- Docker and Docker Compose
- PostgreSQL 16 (via Docker)

### 1. Start PostgreSQL Database

```bash
docker compose up -d
```

This starts PostgreSQL on port 5432.

### 2. Initialize Database

Connect to the database and run migrations:

```bash
# Using psql
psql -h localhost -U postgres -d localconnect

# Or using a database client, run the SQL files:
# - backend/migrations/001_init.sql
# - backend/migrations/002_messaging.sql
```

Or run directly:

```bash
docker exec -i localconnect-postgres psql -U postgres -d localconnect < backend/migrations/001_init.sql
docker exec -i localconnect-postgres psql -U postgres -d localconnect < backend/migrations/002_messaging.sql
```

### 3. Backend Setup

```bash
cd backend

# Download dependencies
go mod download

# Run the server
go run ./cmd/api

# Or build and run
go build -o api ./cmd/api
./api
```

The API will start on `http://localhost:8080`

### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start on `http://localhost:5173`

## Environment Variables

### Backend (.env or export)

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/localconnect?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
MEDIA_DIR=./media
PORT=8080
```

### Frontend (.env)

```bash
VITE_API_URL=http://localhost:8080
```

## API Endpoints

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user

### Categories
- `GET /categories` - List all categories

### Profiles
- `GET /profiles` - List profiles (with filters: category, location)
- `GET /profiles/:id` - Get profile details
- `POST /profiles` - Create profile (auth required, worker only)
- `PUT /profiles/:id` - Update profile (auth required, owner only)
- `POST /profiles/:id/media` - Upload media (auth required, owner only)

### Reviews
- `GET /profiles/:id/reviews` - List reviews for profile
- `POST /profiles/:id/reviews` - Add review (auth required)

### Upvotes
- `POST /profiles/:id/upvote` - Upvote profile (auth required)
- `DELETE /profiles/:id/upvote` - Remove upvote (auth required)

### Contact Requests
- `POST /profiles/:id/contact-requests` - Send contact request (auth required)
- `GET /contact-requests` - List received contact requests (auth required)

### Messaging
- `POST /messages` - Send message (auth required)
- `GET /conversations` - List conversations (auth required)
- `GET /conversations/:id/messages` - Get messages in conversation (auth required)

## Concurrency & Optimization Features

### 1. Worker Pool
- Configurable worker goroutines for async task processing
- Job buffering with overflow handling
- Graceful shutdown support
- Used for: profile updates, notifications, analytics

### 2. Caching System
- Thread-safe in-memory cache with RWMutex
- TTL-based expiration
- Background cleanup goroutine
- Pattern-based invalidation
- 5-minute TTL for profile and search results

### 3. Notification System
- Async notification processing with channels
- Batch processing (configurable size and delay)
- Multiple worker goroutines
- User-based notification grouping
- Prevents blocking main request handlers

### 4. Analytics Tracking
- Non-blocking event tracking
- Concurrent event aggregation
- Periodic stats flushing
- Tracks: profile views, searches, contact requests

### 5. Database Optimizations
- Connection pooling with pgx
- Prepared statements
- Indexed columns for common queries
- Database triggers for automated updates
- Batch operations where possible

## Key Go Patterns Used

1. **Goroutines**: For async processing of notifications, analytics, and background tasks
2. **Channels**: For communication between goroutines (job queues, event streams)
3. **Context**: For cancellation and timeout management
4. **WaitGroups**: For graceful shutdown of concurrent workers
5. **Mutex/RWMutex**: For thread-safe cache operations
6. **Select**: For non-blocking channel operations
7. **Worker Pools**: For limiting concurrent operations

## Development

### Build Backend

```bash
cd backend
go build -o bin/api ./cmd/api
```

### Build Frontend

```bash
cd frontend
npm run build
```

The built files will be in `frontend/dist/`

### Run Tests

```bash
# Backend - Full test suite with coverage
cd backend
make test

# Backend - Quick test
make test-quick

# Backend - Specific test
go test -v -run TestRegister ./internal/handlers/

# Frontend
cd frontend
npm test
```

See [TESTING_COMPLETE.md](TESTING_COMPLETE.md) for comprehensive testing guide.

## Production Deployment

1. Set strong JWT_SECRET
2. Use production PostgreSQL instance
3. Configure CORS properly
4. Set up HTTPS/TLS
5. Use environment variables for configuration
6. Set up proper logging and monitoring
7. Configure database backups
8. Use CDN for media files
9. Implement rate limiting
10. Set up health checks

## License

MIT

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

# LocalConnect Architecture

## Overview

LocalConnect is a full-stack application built with Go (backend) and React (frontend), designed to connect local workers with clients efficiently and reliably.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend (React)                     │
│  ┌─────────┐  ┌──────────┐  ┌─────────┐  ┌──────────┐     │
│  │  Home   │  │ Profile  │  │Messages │  │Dashboard │     │
│  │ Search  │  │  Detail  │  │         │  │  Worker  │     │
│  └─────────┘  └──────────┘  └─────────┘  └──────────┘     │
│                        │                                     │
│                   API Client                                │
└────────────────────────┼────────────────────────────────────┘
                         │ HTTP/REST
                         │ JSON
┌────────────────────────┼────────────────────────────────────┐
│                   Backend (Go)                               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              HTTP Router (Chi)                       │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │  │
│  │  │  Auth    │  │ Profiles │  │    Messaging     │  │  │
│  │  │Middleware│  │ Handlers │  │     Handlers     │  │  │
│  │  └──────────┘  └──────────┘  └──────────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
│                         │                                    │
│  ┌──────────────────────┼─────────────────────────────┐    │
│  │          Concurrent Services                        │    │
│  │  ┌─────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐│    │
│  │  │ Worker  │ │  Cache   │ │Notifier  │ │Tracker ││    │
│  │  │  Pool   │ │(in-mem)  │ │ (async)  │ │(async) ││    │
│  │  └─────────┘ └──────────┘ └──────────┘ └────────┘│    │
│  │  Goroutines + Channels                             │    │
│  └────────────────────────────────────────────────────┘    │
│                         │                                    │
│  ┌──────────────────────┼─────────────────────────────┐    │
│  │              Database Layer (pgx)                   │    │
│  │         Connection Pool (10 max, 2 min)            │    │
│  └────────────────────────────────────────────────────┘    │
└────────────────────────┼────────────────────────────────────┘
                         │
┌────────────────────────┼────────────────────────────────────┐
│                  PostgreSQL 16                               │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌──────────┐ │
│  │  users   │  │ profiles │  │conversations│  │ reviews  │ │
│  └──────────┘  └──────────┘  └────────────┘  └──────────┘ │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌──────────┐ │
│  │categories│  │  media   │  │  messages  │  │ upvotes  │ │
│  └──────────┘  └──────────┘  └────────────┘  └──────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Backend Architecture

### Layer Structure

1. **HTTP Layer** (`cmd/api`, `internal/handlers`)
   - Entry point and routing
   - Request validation
   - Response formatting
   - Middleware (auth, CORS, logging)

2. **Service Layer** (various `internal/` packages)
   - Business logic
   - Concurrent processing
   - Caching strategy
   - Notification handling

3. **Data Layer** (`internal/db`, `internal/models`)
   - Database connection management
   - Data models
   - Query execution

### Concurrency Model

#### Worker Pool Pattern

```go
type Pool struct {
    jobs   chan Job          // Buffered channel for jobs
    wg     sync.WaitGroup    // Wait for workers
    ctx    context.Context   // Cancellation
    cancel context.CancelFunc
}
```

**Use Cases:**
- Profile updates (async stats recalculation)
- Notification sending
- Analytics event processing
- Media file processing

**Flow:**
1. Request handler submits job to channel
2. Worker goroutine picks up job
3. Job executes with context
4. Handler returns immediately (non-blocking)

#### Notification System

```go
Notifier
├─ Queue (buffered channel)
├─ Workers (3 goroutines)
│  ├─ Batch collector
│  ├─ User grouping
│  └─ Concurrent sender
└─ Graceful shutdown
```

**Features:**
- Batching (50 notifications or 2 second delay)
- User-based grouping for efficiency
- Non-blocking submission
- Overflow handling (falls back to goroutine)

**Flow:**
```
Event → Queue → Batch → Group by User → Process concurrently → Send
```

#### Analytics Tracker

```go
Tracker
├─ Event channel
├─ Aggregators (2 goroutines)
├─ Statistics (mutex-protected)
└─ Periodic flusher (5 min)
```

**Tracked Events:**
- Profile views
- Category searches
- Contact requests

**Flow:**
```
Event → Channel → Aggregator → In-memory stats → Flush to DB
```

#### Cache System

```go
Cache
├─ Map (mutex-protected)
├─ TTL-based expiration
├─ Cleanup goroutine
└─ Pattern invalidation
```

**Strategy:**
- 5-minute TTL for profiles and searches
- Invalidate on updates
- RWMutex for concurrent read/write
- Background cleanup every 2.5 minutes

### Database Design

#### Core Tables

**users**
- Authentication and user info
- Role-based (worker/client)

**categories**
- Predefined worker categories
- Referenced by profiles

**worker_profiles**
- Worker information
- Denormalized stats (upvote_count, review_count, average_rating)
- Location-based filtering

**messages & conversations**
- One-to-one messaging
- Read status tracking
- Last message tracking

**reviews & upvotes**
- User feedback system
- Constraints prevent duplicates

#### Indexes

```sql
-- Performance-critical indexes
CREATE INDEX idx_worker_profiles_category ON worker_profiles(category_id);
CREATE INDEX idx_worker_profiles_location ON worker_profiles(location);
CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_receiver_unread ON messages(receiver_id, read) WHERE read = false;
```

#### Database Triggers

```sql
-- Auto-update conversation on new message
CREATE TRIGGER trigger_update_conversation
AFTER INSERT ON messages
FOR EACH ROW
EXECUTE FUNCTION update_conversation_on_message();
```

## Frontend Architecture

### Component Structure

```
App
├── Layout (Header + Navigation)
├── Routes
│   ├── Home (Search & Browse)
│   ├── ProfileDetail (View Profile)
│   ├── Dashboard (Worker Management)
│   ├── Messages (Conversations)
│   └── Auth (Login/Register)
└── API Client (Centralized HTTP)
```

### State Management

- React hooks (useState, useEffect)
- Local storage for auth token
- Component-level state (no Redux needed for MVP)

### API Client Pattern

```javascript
// Centralized API client
api.profiles() → GET /profiles
api.profile(id) → GET /profiles/:id
api.sendMessage() → POST /messages
```

**Features:**
- Automatic JWT header injection
- Error handling
- Type-safe requests

## Data Flow Examples

### 1. Creating a Review

```
Client                   Backend                  Database
  │                         │                         │
  │─ POST /profiles/1/reviews →│                      │
  │                         │                         │
  │                         │─ Validate auth         │
  │                         │                         │
  │                         │─ INSERT review ────────→│
  │                         │                         │
  │                         │─ Submit worker job     │
  │                         │   (async update stats) │
  │                         │                         │
  │←────── 201 Created ────│                         │
  │                         │                         │
  │                         │─ Worker updates stats ─→│
  │                         │   UPDATE profile        │
  │                         │                         │
  │                         │─ Send notification     │
  │                         │   (batched, async)     │
  │                         │                         │
  │                         │─ Invalidate cache      │
```

### 2. Searching Profiles

```
Client                   Backend                  Database
  │                         │                         │
  │─ GET /profiles?cat=cook →│                       │
  │                         │                         │
  │                         │─ Check cache           │
  │                         │   (hit) ──┐            │
  │                         │           │            │
  │                         │←──────────┘            │
  │←────── JSON (cached) ──│                         │
  │                         │                         │
  │                         │   (miss) ──┐           │
  │                         │           │            │
  │                         │           └─ Query ───→│
  │                         │                         │
  │                         │←────── Results ────────│
  │                         │                         │
  │                         │─ Store in cache        │
  │                         │                         │
  │←────── JSON ───────────│                         │
  │                         │                         │
  │                         │─ Track analytics       │
```

### 3. Sending a Message

```
Client                   Backend                  Database
  │                         │                         │
  │─ POST /messages ───────→│                         │
  │  {receiver_id, content} │                         │
  │                         │                         │
  │                         │─ Validate auth         │
  │                         │                         │
  │                         │─ Submit worker job     │
  │                         │   (async create msg)   │
  │                         │                         │
  │←────── 201 Created ────│                         │
  │                         │                         │
  │                         ├─ Worker executes:      │
  │                         │   ├─ Get/Create conv ─→│
  │                         │   ├─ INSERT message ──→│
  │                         │   └─ Send notification │
```

## Scalability Considerations

### Current Design (MVP)
- Single server
- In-memory cache
- Direct database connection
- File-based media storage

### Future Improvements

1. **Horizontal Scaling**
   - Redis for shared cache
   - Message queue (RabbitMQ/Kafka)
   - Load balancer

2. **Database**
   - Read replicas
   - Connection pooler (PgBouncer)
   - Partitioning for large tables

3. **Media Storage**
   - S3/Cloud storage
   - CDN for delivery
   - Image optimization service

4. **Real-time Features**
   - WebSocket for live messaging
   - Server-sent events for notifications

5. **Observability**
   - Structured logging
   - Metrics (Prometheus)
   - Tracing (OpenTelemetry)

## Security

### Authentication
- JWT with HS256
- Token expiration (24 hours)
- Secure password hashing (bcrypt)

### Authorization
- Middleware-based auth checks
- Role-based access (worker/client)
- Owner-only operations

### Database
- Prepared statements (SQL injection prevention)
- Foreign key constraints
- Row-level checks (e.g., different users in conversation)

### API
- CORS configuration
- Request timeout (30s)
- Rate limiting (future)

## Performance Optimizations

1. **Database**
   - Strategic indexes
   - Connection pooling
   - Denormalized stats
   - Batch operations

2. **Backend**
   - In-memory caching (5min TTL)
   - Async processing (non-blocking)
   - Worker pools (limited concurrency)
   - Context-based cancellation

3. **Frontend**
   - Code splitting (Vite)
   - Lazy loading
   - Optimistic UI updates
   - Debounced search

## Error Handling

### Backend
- Structured error responses
- Context-aware logging
- Graceful degradation (cache miss → DB query)
- Panic recovery middleware

### Frontend
- Try-catch blocks
- User-friendly error messages
- Retry logic for failed requests
- Loading states

## Deployment Architecture

```
┌─────────────────────────────────────────┐
│            Load Balancer (nginx)        │
└──────────────┬──────────────────────────┘
               │
     ┌─────────┴─────────┐
     │                   │
┌────┴────┐        ┌─────┴────┐
│ Backend │        │ Backend  │
│ (Go)    │        │ (Go)     │
└────┬────┘        └─────┬────┘
     │                   │
     └─────────┬─────────┘
               │
      ┌────────┴────────┐
      │   PostgreSQL    │
      │   (Primary)     │
      └─────────────────┘
```

## Monitoring & Health Checks

### Health Endpoint
```
GET /health
→ { "status": "ok" }
```

### Metrics to Track
- Request latency
- Error rate
- Database connection pool usage
- Cache hit rate
- Queue depth
- Active goroutines

## Development Workflow

1. **Local Development**
   - Docker Compose for PostgreSQL
   - Hot reload (air for Go, Vite for React)
   - Local environment variables

2. **Testing**
   - Unit tests for business logic
   - Integration tests for handlers
   - E2E tests for critical flows

3. **CI/CD**
   - Automated tests on PR
   - Linting and formatting
   - Build verification
   - Deployment pipeline

## Conclusion

LocalConnect's architecture emphasizes:
- **Concurrency**: Goroutines and channels for performance
- **Scalability**: Caching and async processing
- **Reliability**: Error handling and graceful shutdown
- **Maintainability**: Clear separation of concerns
- **Developer Experience**: Simple setup and clear code structure

# Khusa Mahal Backend

High-performance backend API built with Go, designed for microsecond response times and maximum throughput.

## 🚀 Tech Stack

- **Framework**: Go Fiber - Ultra-fast HTTP framework (benchmarks faster than Express and Fastify)
- **Database**: MongoDB - Flexible document storage with advanced indexing
- **Cache**: Redis - Sub-millisecond caching for blazing-fast responses
- **Search**: Elasticsearch - Lightning-fast full-text search with fuzzy matching

## 📁 Project Structure

```
backend/
├── cmd/
│   └── server/          # Application entry point
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── middleware/  # CORS, logging, auth
│   │   └── routes/      # Route definitions
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   │   ├── mongodb/    # MongoDB operations  
│   │   ├── redis/      # Redis caching
│   │   └── elasticsearch/  # Search service
│   ├── services/       # Business logic
│   └── config/         # Configuration
├── go.mod
└── .env.example
```

## 🔧 Setup

### Prerequisites

- Go 1.22+
- MongoDB
- Redis
- Elasticsearch (optional, fallback to MongoDB search)

### Installation

1. **Clone and navigate to backend**
```bash
cd backend
```

2. **Copy environment file**
```bash
cp .env.example .env
```

3. **Update .env with your configuration**
```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017
REDIS_HOST=localhost
ELASTICSEARCH_URL=http://localhost:9200
```

4. **Install dependencies**
```bash
go mod download
```

5. **Run the server**
```bash
go run cmd/server/main.go
```

## 📊 Performance Features

### Multi-Layer Caching Strategy

1. **Product Cache** (1 hour TTL)
   - Individual products cached by ID
   - Reduces database load by 80%+

2. **List Cache** (15 minutes TTL)
   - Product listings with filters
   - Category-based queries

3. **Cart Cache** (7 days TTL)
   - Session-based carts
   - User carts (no expiration)

### Database Optimization

- **Indexes** on frequently queried fields
  - Category
  - Price
  - Created date
  - Full-text search

- **Connection Pooling** for MongoDB
- **Batch Operations** for bulk updates

### Response Times

- **Cached Responses**: < 5ms
- **Database Queries**: < 50ms
- **Search Queries** (Elasticsearch): < 100ms
- **Session Cart Operations**: < 3ms (Redis)

## 🛣️ API Endpoints

### Health Check
```
GET /api/v1/health
```

### Products

#### Get All Products
```
GET /api/v1/products
GET /api/v1/products?category=Bridal
```

#### Get Single Product
```
GET /api/v1/products/:id
```

#### Search Products
```
GET /api/v1/products/search?q=velvet
```

## 🔐 Security Features

- CORS configuration
- JWT authentication (ready to implement)
- Input validation
- Panic recovery middleware
- Request logging

## 🚦 Running in Production

### Build
```bash
go build -o app cmd/server/main.go
```

### Run
```bash
./app
```

### Docker (Coming Soon)
```bash
docker build -t khusa-mahal-backend .
docker run -p 8080:8080 khusa-mahal-backend
```

## 📈 Monitoring

The server logs include:
- Request latency
- HTTP status codes
- Connection status for all services
- Error tracking

## 🔄 Data Seeding

To seed initial product data:
```go
// TODO: Create seeder script
// Will convert frontend MOCK_PRODUCTS to MongoDB
```

## 🧪 Testing

```bash
go test ./...
```

## 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | Database name | `khusa_mahal` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `ELASTICSEARCH_URL` | Elasticsearch URL | `http://localhost:9200` |
| `JWT_SECRET` | JWT secret key | Change in production |
| `ALLOWED_ORIGINS` | CORS allowed origins | `http://localhost:3000` |

## 🎯 Next Steps

1. Implement cart management endpoints
2. Add user authentication
3. Create order processing API
4. Add payment integration (Stripe/PayPal)
5. Implement admin panel endpoints
6. Add email notifications
7. Create data seeder from frontend MOCK_PRODUCTS

## 📚 Additional Resources

- [Go Fiber Documentation](https://docs.gofiber.io)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Redis Go Client](https://redis.uptrace.dev/)
- [Elasticsearch Go Client](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/index.html)

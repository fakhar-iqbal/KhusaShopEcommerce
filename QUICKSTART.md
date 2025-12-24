# Quick Start Guide

## Prerequisites

You have three options:

### Option 1: MongoDB Only (Recommended for quick testing)
- Install MongoDB only
- Backend will work without Redis and Elasticsearch (with degraded performance)

### Option 2: MongoDB + Redis (Better performance)
- Install MongoDB and Redis
- Get sub-10ms cached responses

### Option 3: Full Stack (Best performance)
- Install MongoDB, Redis, and Elasticsearch
- Get full search capabilities and maximum performance

## Installation

### Windows

**MongoDB:**
```bash
# Download from: https://www.mongodb.com/try/download/community
# Or use Chocolatey:
choco install mongodb
```

**Redis (Optional):**
```bash
# Download from: https://github.com/tporadowski/redis/releases
# Or use Chocolatey:
choco install redis-64
```

**Elasticsearch (Optional):**
```bash
# Download from: https://www.elastic.co/downloads/elasticsearch
```

## Running the Backend

### 1. Start Services

**MongoDB:**
```bash
# Usually starts automatically, or:
mongod
```

**Redis (if installed):**
```bash
redis-server
```

**Elasticsearch (if installed):**
```bash
# Navigate to Elasticsearch folder and run:
bin\elasticsearch.bat
```

### 2. Setup Backend

```bash
cd backend

# Dependencies are already in go.mod
go mod download

# Seed the database with products
go run cmd/seeder/main.go

# Start the server
go run cmd/server/main.go
```

### 3. Test the API

Open your browser or use curl:

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get all products
curl http://localhost:8080/api/v1/products

# Get single product (use ID from above response)
curl http://localhost:8080/api/v1/products/PRODUCT_ID

# Search products
curl http://localhost:8080/api/v1/products/search?q=velvet
```

## Expected Output

When you start the server, you should see:

```
✅ Connected to MongoDB
✅ Connected to Redis
✅ Connected to Elasticsearch
✅ MongoDB indexes created
🚀 Server starting on port 8080
🌐 Health check: http://localhost:8080/api/v1/health
📦 Products API: http://localhost:8080/api/v1/products
```

**Note:** If Redis or Elasticsearch are not running, you'll see warnings but the server will still work!

## Troubleshooting

### MongoDB Connection Failed
- Make sure MongoDB is running: `mongod`
- Check if port 27017 is available

### Redis Connection Failed  
- Redis is optional - server will work without it
- To use Redis: Make sure it's running on port 6379

### Elasticsearch Connection Failed
- Elasticsearch is optional - search will fall back to MongoDB
- To use Elasticsearch: Make sure it's running on port 9200

## Next Steps

Once the backend is running:

1. Test all API endpoints
2. Check MongoDB Compass to see the seeded data
3. Frontend integration (replace MOCK_PRODUCTS with API calls)
4. Add more endpoints (cart, authentication, orders)

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig
	MongoDB       MongoDBConfig
	Redis         RedisConfig
	Elasticsearch ElasticsearchConfig
	JWT           JWTConfig
	CORS          CORSConfig
	Cache         CacheConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type ElasticsearchConfig struct {
	URL   string
	Index string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins string
}

type CacheConfig struct {
	ProductTTL time.Duration
	ListTTL    time.Duration
	CartTTL    time.Duration
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist in production
		fmt.Println("No .env file found, using environment variables")
	}

	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "khusa_mahal"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Elasticsearch: ElasticsearchConfig{
			URL:   getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
			Index: getEnv("ELASTICSEARCH_INDEX", "products"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-this-secret"),
			Expiry: jwtExpiry,
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		},
		Cache: CacheConfig{
			ProductTTL: parseDuration(getEnv("CACHE_PRODUCT_TTL", "3600")),
			ListTTL:    parseDuration(getEnv("CACHE_LIST_TTL", "900")),
			CartTTL:    parseDuration(getEnv("CACHE_CART_TTL", "604800")),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(seconds string) time.Duration {
	s, err := strconv.Atoi(seconds)
	if err != nil {
		return time.Hour
	}
	return time.Duration(s) * time.Second
}

package config

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    ServerPort string
    GinMode    string

    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string

    RedisHost     string
    RedisPort     string
    RedisPassword string
    RedisDB       int

    JWTSecret      string
    JWTExpireHours int
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
    jwtExpireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))

    return &Config{
        ServerPort: getEnv("SERVER_PORT", "8080"),
        GinMode:    getEnv("GIN_MODE", "release"),

        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "studyspot"),
        DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

        RedisHost:     getEnv("REDIS_HOST", "localhost"),
        RedisPort:     getEnv("REDIS_PORT", "6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        RedisDB:       redisDB,

        JWTSecret:      getEnv("JWT_SECRET", "default-secret-key"),
        JWTExpireHours: jwtExpireHours,
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
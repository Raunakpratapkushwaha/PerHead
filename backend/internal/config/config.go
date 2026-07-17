package config

import (
    "log"
    "os"
    "strconv"
    "time"
    "github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	Env                  string
	DBSource             string
	RedisAddr            string
	JWTAccessSecret      string
	JWTRefreshSecret     string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Argon2Memory         uint32
	Argon2Iterations     uint32
	Argon2Parallelism    uint8
	Argon2SaltLength     uint32
	Argon2KeyLength      uint32
}

func LoadConfig(path string) (*Config, error) {
	_ = godotenv.Load(path + "/.env") // Fail silently if file is not found (e.g. in Docker)

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		Env:                  getEnv("ENV", "development"),
		DBSource:             getEnv("DB_SOURCE", ""),
		RedisAddr:            getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		JWTAccessSecret:      getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:     getEnv("JWT_REFRESH_SECRET", ""),
		AccessTokenDuration:  getDurationEnv("ACCESS_TOKEN_DURATION", 15*time.Minute),
		RefreshTokenDuration: getDurationEnv("REFRESH_TOKEN_DURATION", 168*time.Hour),
		Argon2Memory:         uint32(getIntEnv("ARGON2_MEMORY", 65536)),
		Argon2Iterations:     uint32(getIntEnv("ARGON2_ITERATIONS", 1)),
		Argon2Parallelism:    uint8(getIntEnv("ARGON2_PARALLELISM", 2)),
		Argon2SaltLength:     uint32(getIntEnv("ARGON2_SALT_LENGTH", 16)),
		Argon2KeyLength:      uint32(getIntEnv("ARGON2_KEY_LENGTH", 32)),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Warning: failed to parse %s as int: %v", key, err)
		return defaultValue
	}
	return val
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		log.Printf("Warning: failed to parse %s as duration: %v", key, err)
		return defaultValue
	}
	return val
}
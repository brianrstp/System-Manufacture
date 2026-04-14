package config

import "os"

type Config struct {
	ServerPort          string
	MySQLUser           string
	MySQLPassword       string
	MySQLHost           string
	MySQLPort           string
	MySQLDatabase       string
	AdminUser           string
	AdminPassword       string
	AdminJWTSecret      string
	AdminTokenExpiry    string
	CustomerJWTSecret   string
	CustomerTokenExpiry string
}

func Load() Config {
	return Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		MySQLUser:           getEnv("MYSQL_USER", "root"),
		MySQLPassword:       getEnv("MYSQL_PASSWORD", ""),
		MySQLHost:           getEnv("MYSQL_HOST", "127.0.0.1"),
		MySQLPort:           getEnv("MYSQL_PORT", "3306"),
		MySQLDatabase:       getEnv("MYSQL_DATABASE", "manufacture"),
		AdminUser:           getEnv("ADMIN_USER", "admin"),
		AdminPassword:       getEnv("ADMIN_PASSWORD", "admin123"),
		AdminJWTSecret:      getEnv("ADMIN_JWT_SECRET", "supersecretkey"),
		AdminTokenExpiry:    getEnv("ADMIN_TOKEN_EXPIRY", "1h"),
		CustomerJWTSecret:   getEnv("CUSTOMER_JWT_SECRET", "customersecretkey"),
		CustomerTokenExpiry: getEnv("CUSTOMER_TOKEN_EXPIRY", "24h"),
	}
}

func (c Config) DSN() string {
	return c.MySQLUser + ":" + c.MySQLPassword + "@tcp(" + c.MySQLHost + ":" + c.MySQLPort + ")/" + c.MySQLDatabase + "?parseTime=true"
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

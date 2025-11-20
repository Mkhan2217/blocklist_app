package config

import "os"

// getEnv returns the value of the environment variable key or defaultVal if not set
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// Database configuration
func DBUser() string     { return getEnv("CHECKGUARD_DB_USER", "postgres") }
func DBPassword() string { return getEnv("CHECKGUARD_DB_PASSWORD", "root") }
func DBName() string     { return getEnv("CHECKGUARD_DB_NAME", "blocklistdb") }
func DBHost() string     { return getEnv("CHECKGUARD_DB_HOST", "localhost") }
func DBPort() string     { return getEnv("CHECKGUARD_DB_PORT", "5432") }

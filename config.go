package main

import "os"

func DBUser() string {
	if v := os.Getenv("CHECKGUARD_DB_USER"); v != "" {
		return v
	}
	return "postgres"
}
func DBPassword() string {
	if v := os.Getenv("CHECKGUARD_DB_PASSWORD"); v != "" {
		return v
	}
	return "rizzu"
}
func DBName() string {
	if v := os.Getenv("CHECKGUARD_DB_NAME"); v != "" {
		return v
	}
	return "blocklistdb"
}
func DBHost() string {
	if v := os.Getenv("CHECKGUARD_DB_HOST"); v != "" {
		return v
	}
	return "localhost"
}
func DBPort() string {
	if v := os.Getenv("CHECKGUARD_DB_PORT"); v != "" {
		return v
	}
	return "5432"
}

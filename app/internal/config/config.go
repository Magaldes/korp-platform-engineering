package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Database struct{ Host, Port, Name, User, Password string }

func DatabaseFromEnv() (Database, error) {
	db := Database{Host: os.Getenv("KORP_DB_HOST"), Port: os.Getenv("KORP_DB_PORT"), Name: os.Getenv("POSTGRES_DB"), User: os.Getenv("POSTGRES_USER"), Password: os.Getenv("POSTGRES_PASSWORD")}
	if db.Host == "" || db.Port == "" || db.Name == "" || db.User == "" || db.Password == "" {
		return db, fmt.Errorf("database configuration is incomplete")
	}
	return db, nil
}

type Cache struct {
	Address string
	TTL     time.Duration
}

func CacheFromEnv() (Cache, error) {
	address := os.Getenv("KORP_REDIS_ADDR")
	if address == "" {
		address = "redis:6379"
	}
	rawTTL := os.Getenv("KORP_CACHE_TTL_SECONDS")
	if rawTTL == "" {
		rawTTL = "60"
	}
	seconds, err := strconv.Atoi(rawTTL)
	if err != nil || seconds <= 0 {
		return Cache{}, fmt.Errorf("KORP_CACHE_TTL_SECONDS must be a positive integer")
	}
	return Cache{Address: address, TTL: time.Duration(seconds) * time.Second}, nil
}

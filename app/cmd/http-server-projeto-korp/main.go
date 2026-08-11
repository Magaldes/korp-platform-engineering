package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projeto-korp/app/internal/audit"
	"github.com/projeto-korp/app/internal/config"
	"github.com/projeto-korp/app/internal/httpserver"
	"github.com/projeto-korp/app/internal/metrics"
	"github.com/projeto-korp/app/internal/postgres"
	"github.com/projeto-korp/app/internal/redisstats"
)

const listenAddress = ":8080"

func main() {
	metricsComponent := metrics.New()
	cacheConfig, err := config.CacheFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	statisticsCache := redisstats.New(cacheConfig.Address)
	defer statisticsCache.Close()
	var auditService *audit.Service
	if database, err := config.DatabaseFromEnv(); err != nil {
		log.Printf("postgres audit disabled: %v", err)
	} else if pool, err := newPool(database); err != nil {
		log.Printf("postgres audit disabled: %v", err)
	} else {
		repository := postgres.NewRepository(pool)
		defer repository.Close()
		auditService = audit.NewStatisticsService(repository, statisticsCache, cacheConfig.TTL, metricsComponent)
		auditService.StartRetention(context.Background())
	}
	server := &http.Server{
		Addr:    listenAddress,
		Handler: metricsComponent.Handler(httpserver.NewHandler(auditService)),
	}

	log.Printf("%s listening on %s", httpserver.ServiceName, listenAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func newPool(database config.Database) (*pgxpool.Pool, error) {
	connection := url.URL{Scheme: "postgres", User: url.UserPassword(database.User, database.Password), Host: database.Host + ":" + database.Port, Path: database.Name}
	poolConfig, err := pgxpool.ParseConfig(connection.String())
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 0
	poolConfig.MaxConnLifetime = 30 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

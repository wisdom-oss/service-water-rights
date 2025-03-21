package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"microservice/internal"
)

// This file contains the connection to the database which is automatically
// initialized on import/app startup

// Pool is not initialized at the app startup and needs to be initiatized by
// calling [Connect].
var pool *pgxpool.Pool

// Errors which are returned if the database configuration is not in order.
var (
	ErrNoDatabaseUser          = errors.New("no database user configured")
	ErrNoDatabasePassword      = errors.New("no database password configured")
	ErrNoDatabaseHost          = errors.New("no database host configured")
	ErrPoolConfigurationFailed = errors.New("unable to initialize database pool")
	ErrPoolPingFailed          = errors.New("unable to ping database via pool")
)

const (
	KeyUser     = internal.ConfigKey_Postgres_User
	KeyPassword = internal.ConfigKey_Postgres_Password
	KeyHost     = internal.ConfigKey_Postgres_Host
	KeyPort     = internal.ConfigKey_Postgres_Port
	KeySSLMode  = internal.ConfigKey_Postgres_SSLMode
	KeyDatabase = internal.ConfigKey_Postgres_Database
)

const pgSqlConnString = `user=%s password=%s host=%s port=%d sslmode=%s database=%s`

func Pool() *pgxpool.Pool {
	return pool
}

func Connect() (err error) {
	slog.Info("initializing database connection")

	config := internal.Configuration()

	if !config.IsSet("postgres.host") {
		return ErrNoDatabaseHost
	}
	if !config.IsSet("postgres.user") {
		return ErrNoDatabaseUser
	}
	if !config.IsSet("postgres.password") {
		return ErrNoDatabasePassword
	}

	connectionString := fmt.Sprintf(pgSqlConnString,
		config.GetString(KeyUser), config.GetString(KeyPassword),
		config.GetString(KeyHost), config.GetInt(KeyPort),
		config.GetString(KeySSLMode), config.GetString(KeyDatabase),
	)
	slog.Debug("generated connection string", "connString", connectionString)

	slog.Debug("initializing database pool with connection string", "connString", connectionString)
	pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrPoolConfigurationFailed.Error(), err)
	}

	slog.Info("validating database connection")
	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("%s: %w", ErrPoolPingFailed.Error(), err)
	}
	return nil
}

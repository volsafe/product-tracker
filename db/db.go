package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var (
	ErrAlreadyInTX      = errors.New("storage already running in a tx")
	ErrNoTXProvided     = errors.New("no tx provided")
	ErrDBNoTInitiated   = errors.New("db not initiated")
	ErrDBNoRowsEffected = errors.New("db no rows effected")
	ErrMustBeInTx       = errors.New("must be in tx")
	ErrInvalidConfig    = errors.New("invalid database configuration")
)

// DBConfig holds database configuration
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
	SSLMode  string
}

// DB wraps sql.DB with additional functionality
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection with the given parameters
func NewDB(dsn string, maxConns, maxIdle int, timeout time.Duration) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(timeout)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// validateConfig checks if the database configuration is valid
func validateConfig(c struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
}) error {
	if c.User == "" || c.Password == "" || c.Host == "" || c.Port == "" || c.DbName == "" {
		return ErrInvalidConfig
	}
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB == nil {
		return nil
	}
	return db.DB.Close()
}

// Ping checks if the database is accessible
func (d *DB) Ping(ctx context.Context) error {
	if d.DB == nil {
		return ErrDBNoTInitiated
	}
	return d.DB.PingContext(ctx)
}

// BeginTx starts a transaction
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if db.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return db.DB.BeginTx(ctx, opts)
}

// ExecContext executes a query without returning any rows
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return db.DB.ExecContext(ctx, query, args...)
}

// QueryContext executes a query that returns rows
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if db.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return db.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that returns at most one row
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if db.DB == nil {
		return nil
	}
	return db.DB.QueryRowContext(ctx, query, args...)
}

// Stats returns database statistics
func (db *DB) Stats() sql.DBStats {
	if db.DB == nil {
		return sql.DBStats{}
	}
	return db.DB.Stats()
}

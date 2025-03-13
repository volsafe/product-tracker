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

// DBConfig represents the database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

// ValidateConfig validates the database configuration
func ValidateConfig(cfg *DBConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if cfg.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if cfg.User == "" {
		return fmt.Errorf("database user is required")
	}
	if cfg.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if cfg.DbName == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

// DB wraps sql.DB with additional functionality
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(cfg *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
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

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"product-tracker/config"
	"product-tracker/db"

	_ "github.com/lib/pq"
)

// Table and column constants
const (
	tableName = "product_tracker"
	columns   = "name, quantity, energy_consumed, date"
)

// Product represents a product record in the database
type Product struct {
	Name           string  `json:"name" validate:"required"`
	Quantity       int     `json:"quantity" validate:"required,min=0"`
	EnergyConsumed float64 `json:"energy_consumed" validate:"required,min=0"`
	Date           string  `json:"date" validate:"required,datetime=2006-01-02"`
}

// Storage handles database operations for products
type Storage struct {
	db *db.DB
}

// NewStorage creates a new storage instance with a database connection
func NewStorage() (*Storage, error) {
	cfg := config.GetConfig()
	dbConn, err := db.NewDB(cfg.GetDSN(), cfg.Database.MaxConns, cfg.Database.MaxIdle, cfg.Database.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	return &Storage{db: dbConn}, nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

// InsertProduct inserts a new product into the database
func (s *Storage) InsertProduct(ctx context.Context, p Product) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4)", tableName, columns)

	result, err := s.db.ExecContext(ctx, query, p.Name, p.Quantity, p.EnergyConsumed, p.Date)
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return db.ErrDBNoRowsEffected
	}

	return nil
}

// GetProductsByName retrieves products by their name
func (s *Storage) GetProductsByName(ctx context.Context, name string) ([]Product, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE name = $1", columns, tableName)

	rows, err := s.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	return scanProducts(rows)
}

// GetProducts retrieves all products from the database
func (s *Storage) GetProducts(ctx context.Context) ([]Product, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", columns, tableName)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	return scanProducts(rows)
}

// scanProducts scans rows into Product structs
func scanProducts(rows *sql.Rows) ([]Product, error) {
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Name, &p.Quantity, &p.EnergyConsumed, &p.Date); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return products, nil
}

// InsertProducts inserts multiple products in a transaction
func (s *Storage) InsertProducts(ctx context.Context, products []Product) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4)", tableName, columns)

	for _, p := range products {
		result, err := tx.ExecContext(ctx, query, p.Name, p.Quantity, p.EnergyConsumed, p.Date)
		if err != nil {
			return fmt.Errorf("failed to insert product: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return db.ErrDBNoRowsEffected
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetProductsByDateRange retrieves products within a date range
func (s *Storage) GetProductsByDateRange(ctx context.Context, startDate, endDate string) ([]Product, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE date BETWEEN $1 AND $2 ORDER BY date", columns, tableName)

	rows, err := s.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	return scanProducts(rows)
}

// GetProductStats retrieves statistics about products
func (s *Storage) GetProductStats(ctx context.Context) (map[string]interface{}, error) {
	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_products,
			SUM(quantity) as total_quantity,
			SUM(energy_consumed) as total_energy,
			AVG(energy_consumed) as avg_energy
		FROM %s`, tableName)

	var stats struct {
		TotalProducts int     `db:"total_products"`
		TotalQuantity int     `db:"total_quantity"`
		TotalEnergy   float64 `db:"total_energy"`
		AverageEnergy float64 `db:"avg_energy"`
	}

	err := s.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalProducts,
		&stats.TotalQuantity,
		&stats.TotalEnergy,
		&stats.AverageEnergy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get product stats: %w", err)
	}

	return map[string]interface{}{
		"total_products": stats.TotalProducts,
		"total_quantity": stats.TotalQuantity,
		"total_energy":   stats.TotalEnergy,
		"avg_energy":     stats.AverageEnergy,
	}, nil
}

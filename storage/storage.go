package storage

import (
    "context"
    "fmt"
    "product-tracker/db"
    _ "github.com/lib/pq"
)

const columns = "name, quantity, energy_consumed, date"

type Product struct {
    Name           string  `json:"name"`
    Quantity       int     `json:"quantity"`
    EnergyConsumed float64 `json:"energy_consumed"`
    Date           string  `json:"date"`
}

type Storage struct {
    db *db.DB
}

func (s *Storage) Close() error {
    if s.db != nil {
        s.db.Close()
        return nil
    }
    return nil
}

func NewStorage() (*Storage, error) {
    dbConn, err := db.NewDB()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to the database: %w", err)
    }
    return &Storage{db: dbConn}, nil
}

func (s *Storage) InsertProduct(ctx context.Context, p Product) error {
    query := `INSERT INTO product_tracker (` + columns + `) VALUES ($1, $2, $3, $4)`
    stmt, err := s.db.DB.PrepareContext(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to prepare insert statement: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.ExecContext(ctx, p.Name, p.Quantity, p.EnergyConsumed, p.Date)
    if err != nil {
        return fmt.Errorf("failed to execute insert statement: %w", err)
    }

    return nil
}

func (s *Storage) GetProductsByName(ctx context.Context, name string) ([]Product, error) {
    query := "SELECT " + columns + " FROM product_tracker WHERE name = $1"
    stmt, err := s.db.DB.PrepareContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to prepare select statement: %w", err)
    }
    defer stmt.Close()

    rows, err := stmt.QueryContext(ctx, name)
    if err != nil {
        return nil, fmt.Errorf("failed to execute select statement: %w", err)
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.Name, &p.Quantity, &p.EnergyConsumed, &p.Date); err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        products = append(products, p)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }

    return products, nil
}

func (s *Storage) GetProducts(ctx context.Context) ([]Product, error) {
	query := "SELECT " + columns + " FROM product_tracker"
	stmt, err := s.db.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare select statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select statement: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Name, &p.Quantity, &p.EnergyConsumed, &p.Date); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return products, nil
}
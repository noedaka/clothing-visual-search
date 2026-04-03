package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Add(ctx context.Context, product *model.Product) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("failed to rollback the transaction: %v", err)
			}
		}
	}()

	var id int64
	err = tx.QueryRowContext(ctx,
		`INSERT INTO products 
		(name, description, price, category_id)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		product.Name, product.Description, product.Price, product.CategoryID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProductRepo) GetByIDs(ctx context.Context, IDs []int64) ([]model.Product, error) {
	if len(IDs) == 0 {
		return []model.Product{}, nil
	}

	placeholders := make([]string, len(IDs))
	args := make([]any, len(IDs))
	for i, id := range IDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT id, name, description, price, category_id
        FROM products
        WHERE id IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoContent
		}
		
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

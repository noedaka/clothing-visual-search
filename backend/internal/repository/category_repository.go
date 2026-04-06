package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/noedaka/clothing-visual-search/backend/internal/model"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}


func (r *CategoryRepo) Add(ctx context.Context, category string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err 
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("failed to rollback the transaction: %v", err)
			}
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO categories (name) VALUES ($1)",
		category,
	)

	if err != nil {
		return err
	}

	return tx.Commit()

}

func (r *CategoryRepo) List(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name FROM categories`,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoContent
		}
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		if err = rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
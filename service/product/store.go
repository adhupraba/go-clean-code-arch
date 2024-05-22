package product

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adhupraba/ecom/types"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db}
}

func scanRowIntoProduct(row pgx.Row) (*types.Product, error) {
	product := new(types.Product)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM products")

	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)

	for rows.Next() {
		p, err := scanRowIntoProduct(rows)

		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec(
		context.Background(),
		`INSERT INTO products (name, description, image, price, quantity) VALUES ($1, $2, $3, $4, $5)`,
		product.Name, product.Description, product.Image, product.Price, product.Quantity,
	)

	if err != nil {
		return err
	}

	return nil
}

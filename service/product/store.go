package product

import (
	"context"
	"fmt"
	"strconv"

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

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM products")

	if err != nil {
		return nil, err
	}

	products := make([]*types.Product, 0)

	for rows.Next() {
		p, err := scanRowIntoProduct(rows)

		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *Store) GetProductByID(id int) (*types.Product, error) {
	row := s.db.QueryRow(context.Background(), "SELECT * FROM products WHERE id = $1", id)
	p, err := scanRowIntoProduct(row)

	if err != nil {
		return nil, err
	}

	if p.ID == 0 {
		return nil, fmt.Errorf("product not found")
	}

	return p, nil
}

func (s *Store) GetProductsByIDs(ids []int) ([]types.Product, error) {
	params := ""
	args := make([]interface{}, len(ids))

	for i, v := range ids {
		params += `$` + strconv.Itoa(i+1) + `,`
		args[i] = v
	}

	params = params[:len(params)-1] // remove last ','

	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (%s)", params)

	rows, err := s.db.Query(
		context.Background(),
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	products := []types.Product{}

	for rows.Next() {
		p, err := scanRowIntoProduct(rows)

		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec(
		context.Background(),
		"UPDATE products SET name = $1, price = $2, image = $3, description = $4, quantity = $5 WHERE id = $6",
		product.Name, product.Price, product.Image, product.Description, product.Quantity, product.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

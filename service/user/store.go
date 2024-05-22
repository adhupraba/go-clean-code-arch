package user

import (
	"context"
	"fmt"

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

func scanRowIntoUser(row pgx.Row) (*types.User, error) {
	user := new(types.User)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1", email)
	u, err := scanRowIntoUser(row)

	if err != nil {
		return nil, err
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1", id)
	u, err := scanRowIntoUser(row)

	if err != nil {
		return nil, err
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(
		context.Background(),
		`INSERT INTO users ("firstName", "lastName", email, password) VALUES ($1, $2, $3, $4)`,
		user.FirstName, user.LastName, user.Email, user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

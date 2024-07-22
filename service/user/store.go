package user

import (
	"database/sql"
	"errors"

	"github.com/QuangNg14/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	return err
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	user := &types.User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, err
}

// func (s *Store) GetUserByID(id int) (*types.User, error) {
// 	row := s.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
// 	user := &types.User{}
// 	//Scan function copies all the values into the user struct
// 	//after this step the user struct will be populated
// 	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, err
// }

func (s *Store) GetUserByID(id int) (*types.User, error) {
	row := s.db.QueryRow("SELECT id, firstName, lastName, email, password, createdAt FROM users WHERE id = ?", id)
	user := &types.User{}
	//Scan function copies all the values into the user struct
	//after this step the user struct will be populated
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

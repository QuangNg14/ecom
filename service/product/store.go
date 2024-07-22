package product

import (
	"database/sql"

	"github.com/QuangNg14/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	// query the database for all products
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	// create a slice to hold all the products
	products := []types.Product{}
	// iterate over the rows
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// return the products
	return products, nil
}

func scanRowsIntoProduct(rows *sql.Rows) (types.Product, error) {
	p := types.Product{}
	err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.CreatedAt, &p.Description, &p.ImageURL)
	return p, err
}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, price, quantity, image) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.Price, product.Quantity, product.ImageURL)
	return err
}

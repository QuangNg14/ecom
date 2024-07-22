package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/QuangNg14/ecom/service/product"
	"github.com/QuangNg14/ecom/service/user"
	"github.com/gorilla/mux"
)

type APIsServer struct {
	addr string
	db   *sql.DB
}

func NewAPIsServer(addr string, db *sql.DB) *APIsServer {
	return &APIsServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIsServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	// if there is a new version -> increment to version 2

	// This is an easy dependency injection pattern in Go
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Register routes for product (new service)
	productStore := product.NewStore(s.db) // this is from the product table in the database
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	log.Println("Starting server on", s.addr)
	return http.ListenAndServe(s.addr, router)
}

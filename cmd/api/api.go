package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/adhupraba/ecom/service/cart"
	"github.com/adhupraba/ecom/service/order"
	"github.com/adhupraba/ecom/service/product"
	"github.com/adhupraba/ecom/service/user"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{
		addr,
		db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	productStore := product.NewStore(s.db)
	orderStore := order.NewStore(s.db)

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}

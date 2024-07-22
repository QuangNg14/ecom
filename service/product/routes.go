package product

import (
	"log"
	"net/http"

	"github.com/QuangNg14/ecom/types"
	"github.com/QuangNg14/ecom/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	log.Println("Registering /products route")
	router.HandleFunc("/products", h.handleGetProducts).Methods("GET")
	router.HandleFunc("/products", h.handleCreateProduct).Methods("POST")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("handleGetProducts invoked")
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	product := types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
		ImageURL:    payload.Image,
	}

	err = h.store.CreateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}

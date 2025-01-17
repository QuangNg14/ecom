// package cart

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/QuangNg14/ecom/service/auth"
// 	"github.com/QuangNg14/ecom/types"
// 	"github.com/QuangNg14/ecom/utils"
// 	"github.com/go-playground/validator"
// 	"github.com/gorilla/mux"
// )

// type Handler struct {
// 	store      types.ProductStore
// 	orderStore types.OrderStore
// 	userStore  types.UserStore
// }

// func NewHandler(
// 	store types.ProductStore,
// 	orderStore types.OrderStore,
// 	userStore types.UserStore,
// ) *Handler {
// 	return &Handler{
// 		store:      store,
// 		orderStore: orderStore,
// 		userStore:  userStore,
// 	}
// }

// func (h *Handler) RegisterRoutes(router *mux.Router) {
// 	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods("POST")
// }

// func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
// 	userID := auth.GetUserIDFromContext(r.Context())

// 	var cart types.CartCheckoutPayload
// 	if err := utils.ParseJSON(r, &cart); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	if err := utils.Validate.Struct(cart); err != nil {
// 		errors := err.(validator.ValidationErrors)
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
// 		return
// 	}

// 	productIds, err := GetCartItemsIDs(cart.Items)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	// get products
// 	products, err := GetProductsByID(productIds)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	orderID, totalPrice, err := CreateOrder(products, cart.Items, userID)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
// 		"total_price": totalPrice,
// 		"order_id":    orderID,
// 	})
// }

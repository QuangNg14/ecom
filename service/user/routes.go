package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/QuangNg14/ecom/config"
	"github.com/QuangNg14/ecom/service/auth"
	"github.com/QuangNg14/ecom/types"
	"github.com/QuangNg14/ecom/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

/*
Details: It initializes a new Handler struct and returns its pointer.
This is a common pattern in Go for creating and initializing structs.
*/
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

/*
When you see func (h *Handler), it means that the function is a method with a receiver of type *Handler
you can use the h variable to access the fields and methods of the Handler struct. Handler instance
*/
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Register routes here

	// router.HandlerFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// handle login
	var payload types.LoginUserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if the user exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		// user exists
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found invalid email or password"))
		return
	}

	// compare the password
	if !auth.CheckPasswordHash(payload.Password, u.Password) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload from request body
	var payload types.RegisterUserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if the user exists
	_, err = h.store.GetUserByEmail(payload.Email)
	if err == nil {
		// user exists
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}
	if err != sql.ErrNoRows {
		// some other error
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// hashing passwords
	hashedPassword, _ := auth.HashPassword(payload.Password)

	// if not, create a new user
	user := types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	}
	err = h.store.CreateUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("user created: %v", user)

	utils.WriteJSON(w, http.StatusCreated, user)
}

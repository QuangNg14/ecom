package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuangNg14/ecom/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	// Test the user service handlers

	// Create a new user store
	// This user store is a mock users table in the actual database
	userStore := &mockUserStore{users: make(map[string]*types.User)}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		// create a payload
		payload := types.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid",
			Password:  "password123",
		}

		marshalled, _ := json.Marshal(payload)

		// Create a new request
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		// create a new router and register the routes
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		// Call the handler
		router.HandleFunc("/register", handler.handleRegister).Methods("POST")
		router.ServeHTTP(rr, req)

		// Check the response
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code to be 400, got %d", rr.Code)
		}
	})

	/*
		Explanation of the Mock and Tests
		Mock Implementation:

		mockUserStore simulates the UserStore interface using an in-memory map.
		Methods GetUserByEmail, GetUserByID, and CreateUser interact with this map instead of a real database.
		Using the Mock in Tests:

		Initialize mockUserStore with an empty map: userStore := &mockUserStore{users: make(map[string]*types.User)}.
		Pass the mock store to the handler: handler := NewHandler(userStore).
		Write test cases to verify the handler's behavior using httptest.NewRecorder to simulate HTTP requests and responses.
	*/
	t.Run("should correctly register the user", func(t *testing.T) {
		// create a payload
		payload := types.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "Hello@gmail.com",
			Password:  "password123",
		}

		marshalled, _ := json.Marshal(payload)

		// Create a new request
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		// create a new router
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		// Call the handler
		router.HandleFunc("/register", handler.handleRegister).Methods("POST")
		router.ServeHTTP(rr, req)

		// Log the response body for debugging
		t.Logf("Response body: %s", rr.Body.String())

		// Check the response
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code to be %d, got %d", http.StatusCreated, rr.Code)
		}
	})

}

type mockUserStore struct {
	users map[string]*types.User
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserStore) CreateUser(user types.User) error {
	if _, exists := m.users[user.Email]; exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}
	m.users[user.Email] = &user
	return nil
}

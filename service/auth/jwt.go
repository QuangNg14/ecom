package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/QuangNg14/ecom/config"
	"github.com/QuangNg14/ecom/types"
	"github.com/QuangNg14/ecom/utils"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the request
		tokenString := utils.GetTokenFromRequest(r)
		log.Printf("Token received: %s", tokenString)

		// validate the JWT token
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// Extract claims from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("failed to parse token claims")
			permissionDenied(w)
			return
		}

		// Get userID from claims
		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			log.Println("userID not found in token claims")
			permissionDenied(w)
			return
		}
		userID := int(userIDFloat)
		log.Printf("UserID extracted: %d", userID)

		// fetch the user from the DB
		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}
		if u == nil {
			log.Println("user is nil")
			permissionDenied(w)
			return
		}
		log.Printf("User retrieved: %+v", u)

		// set context "userID" to the user ID
		// Add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Duration(config.Envs.JWTExpirationInSecond) * time.Second
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    float64(userID), // JWT claims are typically floats
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}
	return token, nil
}

func permissionDenied(w http.ResponseWriter) {
	http.Error(w, "permission denied", http.StatusForbidden)
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}

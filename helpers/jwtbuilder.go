package jwtbuilder

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TicketExpiredError struct {
	message string
}

func (e *TicketExpiredError) Error() string {
	return e.message
}

//CreateJWTToken creates a new token
func CreateJWTToken(userID string) (string, error) {
	// Set the claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the encoded token
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyJWTToken verifies the token
func VerifyJWTToken(tokenString string) (string, error) {

	// Parse the token
	var err error
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return secretKey, nil
	})

	// Validate the token and return the custom claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return "", &TicketExpiredError{message: "Ticket expired"}

		}
		userID := claims["user_id"].(string)
		return userID, nil
	}
	return "", err
}

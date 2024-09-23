package tokenjwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtkey = []byte("SECRET_KEY")

type Claims struct {
	Email string `json:"eamil"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	Claims := &Claims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, err := token.SignedString(jwtkey)
	return tokenString, err
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

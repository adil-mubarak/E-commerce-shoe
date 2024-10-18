package tokenjwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtkey = []byte("SECRET_KEY")
var refreshKey = []byte("REFRESH_SECRET_KEY")

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint, email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	Claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenString, err := token.SignedString(jwtkey)
	return tokenString, err
}

func RefreshJWT(userID uint, email string,role string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := &RefreshClaims{
		UserID: userID,
		Email:  email,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	tokenString, err := refreshToken.SignedString(refreshKey)
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

func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}

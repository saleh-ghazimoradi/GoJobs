package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:"userid"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

func GenerateToken(username string, userID int64, isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(10 * time.Hour)
	claims := &Claims{
		Username: username,
		UserID:   userID,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.AppConfig.JWT.SecretKEY))
}

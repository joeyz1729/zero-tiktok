package jwtx

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	AccessTokenExpireDuration  = 14 * 24 * time.Hour
	RefreshTokenExpireDuration = 30 * 24 * time.Hour

	ErrorInvalidToken = errors.New("invalid token")

	secret = []byte("baldur's gate 3")
	issuer = "zero-tiktok"
)

type Claims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenToken(userId int64, username string) (aToken, rToken string, err error) {
	c := &Claims{
		userId,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpireDuration).Unix(),
			Issuer:    issuer,
		},
	}
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(secret)
	if err != nil {
		return "", "", err
	}
	//rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
	//	ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(),
	//	Issuer:    issuer,
	//}).SignedString(secret)
	//if err != nil {
	//	return aToken, "", err
	//}

	return aToken, rToken, nil
}

func ParseToken(tokenStr string) (claims *Claims, err error) {
	claims = new(Claims)
	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)
	//fmt.Printf("%v, %v\n", token, claims)
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claims, nil
	}
	return nil, ErrorInvalidToken
}

func keyFunc(_ *jwt.Token) (any interface{}, err error) {
	return secret, nil
}

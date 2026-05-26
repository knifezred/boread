package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret []byte
var jwtExpire time.Duration

func Init(secret string, expire int) {
	jwtSecret = []byte(secret)
	jwtExpire = time.Duration(expire) * time.Second
}

func GenerateToken(userID uint64, username string) (token string, expiresAt int64, err error) {
	expiresAtTime := time.Now().Add(jwtExpire)
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAtTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenObj.SignedString(jwtSecret)
	return token, expiresAtTime.Unix(), err
}

// GenerateRefreshToken 生成刷新token，过期时间为token的2倍
func GenerateRefreshToken(userID uint64, username string) (token string, expiresAt int64, err error) {
	expiresAtTime := time.Now().Add(jwtExpire * 2)
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAtTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenObj.SignedString(jwtSecret)
	return token, expiresAtTime.Unix(), err
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

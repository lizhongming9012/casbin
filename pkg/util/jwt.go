package util

import (
	"NULL/casbin/pkg/setting"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

type CustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)

	claims := CustomClaims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "NULL",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (interface{}, error) {
	secret := []byte(setting.AppSetting.JwtSecret)
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected siging method:%v", token.Header["alg"])
		}
		return secret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok {
			if tokenClaims.Valid {
				return claims, nil
			} else {
				if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
					return claims, err
				}
			}
		}
	}
	return nil, err
}

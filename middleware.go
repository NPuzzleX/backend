package main

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func createJWT(claims jwt.MapClaims) (string, error) {
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func parseJWT(jwtString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(jwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func validateToken(token string, refresh bool, needLogin bool) (string, error) {
	if token == "" {
		return "", errors.New("NEEDS TOKEN")
	}

	if len(strings.Split(token, " ")) != 2 {
		return "", errors.New("INVALID TOKEN")
	}

	if strings.Split(token, " ")[0] != "Bearer" {
		return "", errors.New("INVALID TOKEN")
	}

	token = strings.Split(token, " ")[1]

	result, err := parseJWT(token)

	if err != nil {
		return "", err
	} else if result["expire"] == nil {
		return "", errors.New("INVALID TOKEN")
	} else if result["id"] == nil {
		return "", errors.New("INVALID TOKEN")
	} else if needLogin && (result["id"].(string) == "") {
		return "", errors.New("NEED LOGIN")
	} else {
		expireTime, err := time.Parse(time.RFC1123Z, result["expire"].(string))
		if err != nil {
			return "", errors.New("INVALID TOKEN")
		}

		if refresh {
			return result["id"].(string), nil
		} else {
			if expireTime.After(time.Now()) {
				return result["id"].(string), nil
			} else {
				return "", errors.New("TOKEN EXPIRED")
			}
		}
	}
}

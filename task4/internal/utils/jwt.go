package utils

import (
	"context"
	"encoding/json"
	"go-study/task4/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, username string, config config.Config) (string, error) {
	now := time.Now()
	expireTime := now.Add(time.Duration(config.JWT.AccessExpire) * time.Second)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "blog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT.AccessSecret))
}

func ParseJWT(tokenString string, config config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// 从上下文获取用户ID - 修复版本
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	if value := ctx.Value("user_id"); value != nil {
		if number, ok := value.(json.Number); ok {
			if intVal, err := number.Int64(); err == nil {
				return uint(intVal), true
			}
		}
	}
	return 0, false
}

// 从上下文获取用户名 - 修复版本
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	if username, ok := ctx.Value("username").(string); ok {
		return username, true
	}
	return "", false
}

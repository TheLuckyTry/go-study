package utils

import (
	"context"
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
	// 尝试从不同可能的键中获取用户ID
	if userId, ok := ctx.Value("userId").(uint); ok {
		return userId, true
	}
	if userId, ok := ctx.Value("userId").(int64); ok {
		return uint(userId), true
	}
	if userId, ok := ctx.Value("userId").(float64); ok {
		return uint(userId), true
	}
	if userId, ok := ctx.Value("userId").(int); ok {
		return uint(userId), true
	}

	// 尝试从 JWT 标准键中获取
	if userId, ok := ctx.Value("userID").(uint); ok {
		return userId, true
	}
	if userId, ok := ctx.Value("user_id").(uint); ok {
		return userId, true
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

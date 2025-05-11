package main

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 通过用户ID获取JWT字符串
func GetJwt(userId uint, jwtKey string) (string, error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		},
	).SignedString([]byte(jwtKey))
}

// 制作响应体
func Resp(message string, err error, data any) gin.H {
	var errString *string
	if err != nil {
		data := err.Error()
		errString = &data
	}
	return gin.H{
		"message": message,
		"error":   errString,
		"data":    data,
	}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func Select(columns ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(columns)
	}
}

package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/bestcb2333/gin-gorm-preloader/preloader"
)

func GetRouter(db *gorm.DB, config *Config) *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	bc := preloader.BaseConfig{
		DB:            db,
		JWTKey:        config.JWTKey,
		JWTExpHours:   time.Hour * 24,
		UserTableName: "users",
		AdminColName:  "admin",
		Resp: func(message string, err error, data any) gin.H {
			var errStr *string
			if err != nil {
				str := err.Error()
				errStr = &str
			}
			return gin.H{
				"message": message,
				"error":   errStr,
				"data":    data,
			}
		},
	}

	RegisterPatentHandler(r, &bc)
	RegisterReportHandler(r, &bc)
	RegisterKeywordHandler(r, &bc)
	RegisterUserHandler(r, &bc)
	RegisterSuggestionHandler(r, &bc)

	return r
}

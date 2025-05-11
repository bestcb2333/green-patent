package main

import (
	"github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListUserReq struct {
	preloader.PageConfig
}

func RegisterUserHandler(r *gin.Engine, bc *preloader.BaseConfig) {

	r.GET("/users", preloader.CreateListHandler[User](
		&preloader.Config[ListUserReq]{
			Base: bc,
			DefReq: ListUserReq{PageConfig: preloader.PageConfig{
				Page:     1,
				PageSize: 10,
			}},
		},
		func(q *gorm.DB, c *gin.Context, u *User, r *ListUserReq) *gorm.DB {

			return q
		},
	))

	r.GET("/users/:id", preloader.CreateGetHandler[User](
		&preloader.Config[struct{}]{
			Base: bc,
		},
		func(q *gorm.DB, c *gin.Context, u *User) *gorm.DB {
			return q.Preload(
				"Patents", Select("id", "user_id", "name"),
			).Preload(
				"Reports", Select("id", "user_id", "name"),
			).Preload(
				"Suggestions", Select("id", "user_id", "title"),
			).Preload(
				"Keywords", Select("id", "user_id", "value"),
			)
		},
	))
}

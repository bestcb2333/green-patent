package main

import (
	"github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListKeywordReq struct {
	preloader.PageConfig
}

func RegisterKeywordHandler(r *gin.Engine, bc *preloader.BaseConfig) {

	r.GET("/keywords", preloader.CreateListHandler[Keyword](
		&preloader.Config[ListKeywordReq]{
			Base: bc,
			DefReq: ListKeywordReq{PageConfig: preloader.PageConfig{
				Page:     1,
				PageSize: 10,
			}},
		},
		func(q *gorm.DB, c *gin.Context, u *User, r *ListKeywordReq) *gorm.DB {
			return q.Preload("User", Select("id", "nickname"))
		},
	))

	r.GET("/keywords/:id", preloader.CreateGetHandler[Keyword](
		&preloader.Config[struct{}]{
			Base: bc,
		},
		func(q *gorm.DB, c *gin.Context, u *User) *gorm.DB {
			return q.Preload("Patents").Preload("Reports")
		},
	))

}

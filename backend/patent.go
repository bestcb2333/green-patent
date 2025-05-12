package main

import (
	"github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListPatentReq struct {
	preloader.PageConfig
}

func RegisterPatentHandler(r *gin.Engine, bc *preloader.BaseConfig) {

	r.GET("/patents", preloader.CreateListHandler[Patent](
		&preloader.Config[ListPatentReq]{
			Base:   bc,
			DefReq: ListPatentReq{},
		},
		func(q *gorm.DB, c *gin.Context, u *User, r *ListPatentReq) *gorm.DB {
			q = q.Preload(
				"User", Select("id", "nickname"),
			).Preload(
				"Keywords", Select("id", "value"),
			)
			return q
		},
	))

	r.GET("/patents/:id", preloader.CreateGetHandler[Patent](
		&preloader.Config[struct{}]{
			Base: bc,
		},
		func(q *gorm.DB, c *gin.Context, u *User) *gorm.DB {
			return q.Preload("User", Select("id", "nickname")).Preload("Keywords", Select("id", "value"))
		},
	))

	r.GET("/stats/patents", preloader.Preload(
		&preloader.Config[struct{}]{
			Base: bc,
		},
		func(c *gin.Context, u *User, r *struct{}) {

			query := bc.DB.Model(new(Patent))

			var total int64
			if err := query.Count(&total).Error; err != nil {
				c.JSON(500, bc.Resp("统计总数失败", err, nil))
				return
			}

			var solved int64
			if err := query.Where("status = true").Count(&solved).Error; err != nil {
				c.JSON(500, bc.Resp("统计已处理数量失败", err, nil))
				return
			}

			c.JSON(200, bc.Resp("统计成功", nil, gin.H{
				"total":   total,
				"solved":  solved,
				"pending": total - solved,
			}))
		},
	))
}

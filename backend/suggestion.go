package main

import (
	"github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListSuggestionReq struct {
	PatentID uint `form:"patent_id"`
	ReportID uint `form:"report_id"`
	preloader.PageConfig
}

func RegisterSuggestionHandler(r *gin.Engine, bc *preloader.BaseConfig) {

	r.GET("/suggestions", preloader.CreateListHandler[Suggestion](
		&preloader.Config[ListSuggestionReq]{
			Base: bc,
			DefReq: ListSuggestionReq{PageConfig: preloader.PageConfig{
				Page:     1,
				PageSize: 10,
			}},
		},
		func(q *gorm.DB, c *gin.Context, u *User, r *ListSuggestionReq) *gorm.DB {
			q = q.Preload(
				"User", Select("id", "nickname"),
			).Preload(
				"Patent", Select("id", "name"),
			).Preload(
				"Report", Select("id", "name"),
			)
			if r.PatentID != 0 {
				q = q.Where("patent_id = ?", r.PatentID)
			}
			if r.ReportID != 0 {
				q = q.Where("report_id = ?", r.ReportID)
			}
			return q
		},
	))

	r.GET("/suggestions/:id", preloader.CreateGetHandler[Suggestion](
		&preloader.Config[struct{}]{
			Base: bc,
		},
		func(q *gorm.DB, c *gin.Context, u *User) *gorm.DB {
			return q.Preload(
				"User", Select("id", "nickname"),
			).Preload(
				"Patent", Select("id", "name"),
			).Preload(
				"Report", Select("id", "name"),
			)
		},
	))
}

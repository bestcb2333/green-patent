package main

import (
	"github.com/bestcb2333/gin-gorm-preloader/preloader"
	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListReportReq struct {
	preloader.PageConfig
	Status bool `form:"status"`
}

func RegisterReportHandler(r *gin.Engine, bc *preloader.BaseConfig, config *Config) {

	r.GET("/reports", preloader.CreateListHandler[Report](
		&preloader.Config[ListReportReq]{
			Base:   bc,
			DefReq: ListReportReq{},
		},
		func(q *gorm.DB, c *gin.Context, u *User, r *ListReportReq) *gorm.DB {
			q = q.Preload(
				"User", Select("id", "nickname"),
			).Preload(
				"Keywords", Select("id", "value"),
			)
			if r.Status {
				q = q.Where("status = ?", false)
			}
			return q
		},
	))

	r.GET("/reports/:id", preloader.CreateGetHandler[Report](
		&preloader.Config[struct{}]{
			Base: bc,
		},
		func(q *gorm.DB, c *gin.Context, u *User) *gorm.DB {
			return q.Preload("User", Select("id", "nickname")).Preload("Keywords", Select("id", "value"))
		},
	))

	r.GET("/trend", preloader.Preload(
		&preloader.Config[struct{}]{
			Base: bc,
		},
		func(c *gin.Context, u *User, r *struct{}) {

			var data []int64
			years := []uint{2021, 2022, 2023, 2024, 2025}

			for _, year := range years {
				var count int64
				if err := bc.DB.Model(new(Report)).Where("year = ?", year).Count(&count).Error; err != nil {
					c.JSON(500, bc.Resp("计数失败", err, nil))
					return
				}
				data = append(data, count)
			}

			c.JSON(200, bc.Resp("获取成功", nil, data))
		},
	))

	r.GET("/stats/reports", preloader.Preload(
		&preloader.Config[struct{}]{
			Base: bc,
		},
		func(c *gin.Context, u *User, r *struct{}) {

			query := bc.DB.Model(new(Report))

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

	r.POST("/reports", p.Preload(
		&p.Config[ReportDTO]{
			Base: bc,
			Bind: &p.BindConfig{Multipart: true},
		},
		func(c *gin.Context, u *User, r *ReportDTO) {

			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(500, bc.Resp("读取文件失败", err, nil))
				return
			}

			if err := bc.DB.Where("file = ?", file.Filename).First(new(Report)).Error; err == nil {
				c.JSON(400, bc.Resp("已存在同名文件", err, nil))
				return
			}

			if err := c.SaveUploadedFile(file, config.Path.Report+file.Filename); err != nil {
				c.JSON(500, bc.Resp("文件保存失败", err, nil))
				return
			}

			data := new(Report)
			data.ReportDTO = r
			data.File = file.Filename
			data.UserID = 1
			data.Name = r.Name
			if err := bc.DB.Create(data).Error; err != nil {
				c.JSON(500, bc.Resp("数据库记录创建失败", err, nil))
				return
			}

			c.JSON(200, bc.Resp("上传成功", nil, nil))
		},
	))
}

package main

import (
	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListPatentReq struct {
	p.PageConfig
	Status bool `form:"status"`
}

func RegisterPatentHandler(r *gin.Engine, bc *p.BaseConfig, config *Config) {

	r.GET("/patents", p.CreateListHandler[Patent](
		&p.Config[ListPatentReq]{
			Base:   bc,
			DefReq: ListPatentReq{},
		},
		func(q *gorm.DB, c *gin.Context, u *User, r *ListPatentReq) *gorm.DB {
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

	r.GET("/patents/:id", p.CreateGetHandler[Patent](
		&p.Config[struct{}]{
			Base: bc,
		},
		func(q *gorm.DB, c *gin.Context, u *User) *gorm.DB {
			return q.Preload("User", Select("id", "nickname")).Preload("Keywords", Select("id", "value"))
		},
	))

	r.GET("/stats/patents", p.Preload(
		&p.Config[struct{}]{
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

	r.POST("/patents", p.Preload(
		&p.Config[PatentDTO]{
			Base: bc,
			Bind: &p.BindConfig{Multipart: true},
		},
		func(c *gin.Context, u *User, r *PatentDTO) {

			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(500, bc.Resp("读取文件失败", err, nil))
				return
			}

			if err := bc.DB.Where("file = ?", file.Filename).First(new(Patent)).Error; err == nil {
				c.JSON(400, bc.Resp("已存在同名文件", err, nil))
				return
			}

			if err := c.SaveUploadedFile(file, config.Path.Patent+file.Filename); err != nil {
				c.JSON(500, bc.Resp("文件保存失败", err, nil))
				return
			}

			data := new(Patent)
			data.PatentDTO = r
			data.File = file.Filename
			data.UserID = 1
			if err := bc.DB.Create(data).Error; err != nil {
				c.JSON(500, bc.Resp("数据库记录创建失败", err, nil))
				return
			}

			c.JSON(200, bc.Resp("上传成功", nil, nil))
		},
	))
}

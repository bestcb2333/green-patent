package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库
func InitDB(config *Config) (*gorm.DB, error) {

	// 连接数据库
	db, err := gorm.Open(mysql.Open(fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB.User,
		config.DB.Pass,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
	)))
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库表
	for _, value := range []any{
		new(User), new(Patent), new(Report), new(Suggestion), new(Keyword),
	} {
		db.AutoMigrate(value)
	}

	return db, nil
}

type IDField struct {
	ID uint `json:"id" gorm:"primarykey;comment:ID"`
}

type CreatedAtField struct {
	CreatedAt *time.Time `json:"createdAt" gorm:"comment:创建时间"`
}

type User struct {
	IDField
	CreatedAtField
	Name        string       `json:"name" gorm:"type:VARCHAR(20);not null;unique;comment:用户名"`
	Nickname    string       `json:"nickname" gorm:"type:VARCHAR(20);not null;comment:昵称"`
	Password    string       `json:"-" gorm:"type:CHAR(64);not null;comment:密码"`
	Email       string       `json:"email" gorm:"type:VARCHAR(50);not null;unique;comment:邮箱"`
	Phone       string       `json:"phone" gorm:"type:VARCHAR(20);not null;unique;comment:手机号"`
	Admin       bool         `json:"admin" gorm:"not null;comment:是管理员"`
	Patents     []Patent     `json:"patents" gorm:"constraint:OnDelete:CASCADE"`
	Reports     []Report     `json:"reports" gorm:"constraint:OnDelete:CASCADE"`
	Suggestions []Suggestion `json:"suggestions" gorm:"constraint:OnDelete:CASCADE"`
	Keywords    []Keyword    `json:"keywords" gorm:"constraint:OnDelete:CASCADE"`
}

type PatentDTO struct {
	Name   string `json:"name" form:"name" gorm:"type:VARCHAR(50);not null;unique;comment:专利名称"`
	Number string `json:"number" form:"number" gorm:"type:VARCHAR(20);not null;unique;comment:专利号"`
}

type Patent struct {
	IDField
	CreatedAtField
	*PatentDTO
	File     string    `json:"file" gorm:"type:VARCHAR(50);not null;unique;comment:文件名"`
	Status   bool      `json:"status" gorm:"not null;comment:状态"`
	UserID   uint      `json:"-" gorm:"index;not null;comment:上传者"`
	User     *User     `json:"user" gorm:"constraint:OnDelete:CASCADE"`
	Keywords []Keyword `json:"keywords" gorm:"many2many:patent_keywords;constraint:OnDelete:CASCADE"`
}

type ReportDTO struct {
	Name string `json:"name" form:"name" gorm:"type:VARCHAR(20);not null;unique;comment:年报名称"`
	Year uint   `json:"year" form:"year" gorm:"not null;comment:年份"`
}

type Report struct {
	IDField
	CreatedAtField
	*ReportDTO
	Status   bool      `json:"status" gorm:"not null;comment:状态"`
	File     string    `json:"file" gorm:"type:VARCHAR(50);not null;unique;comment:文件名"`
	UserID   uint      `json:"-" gorm:"index;not null;comment:上传者"`
	User     *User     `json:"user" gorm:"constraint:OnDelete:CASCADE"`
	Keywords []Keyword `json:"keywords" gorm:"many2many:report_keywords;constraint:OnDelete:CASCADE"`
}

type SuggestionDTO struct {
	Title   string `json:"title" gorm:"type:VARCHAR(50);not null;unique;comment:标题"`
	Content string `json:"content" gorm:"type:TEXT;not null;comment:内容"`
}

type Suggestion struct {
	IDField
	CreatedAtField
	User     *User   `json:"user" gorm:"constraint:OnDelete:CASCADE"`
	UserID   uint    `json:"-" gorm:"index;not null;comment:创建者ID"`
	PatentID *uint   `json:"-" gorm:"index;comment:专利ID"`
	Patent   *Patent `json:"patent" gorm:"constraint:OnDelete:CASCADE"`
	ReportID *uint   `json:"-" gorm:"index;comment:年报ID"`
	Report   *Report `json:"report" gorm:"constraint:OnDelete:CASCADE"`
	*SuggestionDTO
}

type KeywordDTO struct {
	Value string `json:"value" gorm:"type:VARCHAR(20);not null;comment:关键词"`
}

type Keyword struct {
	IDField
	CreatedAtField
	UserID uint  `json:"-" gorm:"index;not null;comment:创建者"`
	User   *User `json:"user" gorm:"constraint:OnDelete:CASCADE"`
	*KeywordDTO
	Patents []Patent `json:"patents" gorm:"many2many:patent_keywords;constraint:OnDelete:CASCADE"`
	Reports []Report `json:"reports" gorm:"many2many:report_keywords;constraint:OnDelete:CASCADE"`
}

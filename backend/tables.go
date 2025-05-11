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
	ID uint `json:"id,omitempty" gorm:"primarykey;comment:ID"`
}

type CreatedAtField struct {
	CreatedAt *time.Time `json:"createdAt,omitempty" gorm:"comment:创建时间"`
}

type User struct {
	IDField
	CreatedAtField
	Name        string       `json:"name,omitempty" gorm:"type:VARCHAR(20);not null;unique;comment:用户名"`
	Nickname    string       `json:"nickname,omitempty" gorm:"type:VARCHAR(20);not null;comment:昵称"`
	Password    string       `json:"-" gorm:"type:CHAR(64);not null;comment:密码"`
	Email       string       `json:"email,omitempty" gorm:"type:VARCHAR(50);not null;unique;comment:邮箱"`
	Phone       string       `json:"phone,omitempty" gorm:"type:VARCHAR(20);not null;unique;comment:手机号"`
	Admin       bool         `json:"admin,omitempty" gorm:"not null;comment:是管理员"`
	Patents     []Patent     `json:"patents,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Reports     []Report     `json:"reports,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Suggestions []Suggestion `json:"suggestions,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Keywords    []Keyword    `json:"keywords,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

type PatentDTO struct {
	Name   string `json:"name,omitempty" gorm:"type:VARCHAR(50);not null;unique;comment:专利名称"`
	Number string `json:"number,omitempty" gorm:"type:VARCHAR(20);not null;unique;comment:专利号"`
	Status bool   `json:"status,omitempty" gorm:"not null;comment:状态"`
	File   string `json:"file,omitempty" gorm:"type:VARCHAR(50);not null;unique;comment:文件名"`
}

type Patent struct {
	IDField
	CreatedAtField
	*PatentDTO
	UserID   uint      `json:"-" gorm:"index;not null;comment:上传者"`
	User     *User     `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Keywords []Keyword `json:"keywords,omitempty" gorm:"many2many:patent_keywords;constraint:OnDelete:CASCADE"`
}

type ReportDTO struct {
	Name   string `json:"name,omitempty" gorm:"type:VARCHAR(20);not null;unique;comment:年报名称"`
	Year   uint   `json:"year,omitempty" gorm:"not null;comment:年份"`
	Status bool   `json:"status,omitempty" gorm:"not null;comment:状态"`
	File   string `json:"file,omitempty" gorm:"type:VARCHAR(50);not null;unique;comment:文件名"`
}

type Report struct {
	IDField
	CreatedAtField
	*ReportDTO
	UserID   uint      `json:"-" gorm:"index;not null;comment:上传者"`
	User     *User     `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Keywords []Keyword `json:"keywords,omitempty" gorm:"many2many:report_keywords;constraint:OnDelete:CASCADE"`
}

type SuggestionDTO struct {
	Title   string `json:"title,omitempty" gorm:"type:VARCHAR(50);not null;unique;comment:标题"`
	Content string `json:"content,omitempty" gorm:"type:TEXT;not null;comment:内容"`
}

type Suggestion struct {
	IDField
	CreatedAtField
	User     *User   `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	UserID   uint    `json:"-" gorm:"index;not null;comment:创建者ID"`
	PatentID *uint   `json:"-" gorm:"index;comment:专利ID"`
	Patent   *Patent `json:"patent,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	ReportID *uint   `json:"-" gorm:"index;comment:年报ID"`
	Report   *Report `json:"report,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	*SuggestionDTO
}

type KeywordDTO struct {
	Value string `json:"value,omitempty" gorm:"type:VARCHAR(20);not null;comment:关键词"`
}

type Keyword struct {
	IDField
	CreatedAtField
	UserID uint  `json:"-" gorm:"index;not null;comment:创建者"`
	User   *User `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	*KeywordDTO
	Patents []Patent `json:"patents,omitempty" gorm:"many2many:patent_keywords;constraint:OnDelete:CASCADE"`
	Reports []Report `json:"reports,omitempty" gorm:"many2many:report_keywords;constraint:OnDelete:CASCADE"`
}

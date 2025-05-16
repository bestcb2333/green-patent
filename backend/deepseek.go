package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	p "github.com/bestcb2333/gin-gorm-preloader/preloader"
	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type DeepSeekKeywords struct {
	Keywords []uint `json:"keywords"`
}

type ExtractReq struct {
	ID uint `uri:"id"`
}

func RegDeepSeekhandler(r *gin.Engine, bc *p.BaseConfig, cfg *Config) {

	client := openai.NewClient(
		option.WithAPIKey(cfg.DeepSeekAPI),
		option.WithBaseURL("https://api.deepseek.com"),
	)

	r.GET("/extract/patent/:id", p.Preload(
		&p.Config[ExtractReq]{
			Base: bc,
			Bind: &p.BindConfig{Param: true},
		},
		func(c *gin.Context, u *User, r *ExtractReq) {

			var patent Patent
			if err := bc.DB.Where("id = ?", r.ID).First(&patent).Error; err != nil {
				c.JSON(500, bc.Resp("获取专利失败", err, nil))
				return
			}

			content, err := GetPDFContent(filepath.Join(cfg.Path.Patent, patent.File))
			if err != nil {
				c.JSON(500, bc.Resp("提取PDF内容失败", err, nil))
				return
			}

			var allKeywords []Keyword
			if err := bc.DB.Select("id", "value").Find(&allKeywords).Error; err != nil {
				c.JSON(500, bc.Resp("查询失败", err, nil))
				return
			}

			byteKeywords, err := json.Marshal(&allKeywords)
			if err != nil {
				c.JSON(500, bc.Resp("序列化失败", err, nil))
				return
			}

			SystemPrompt := fmt.Sprintf(`这是关键词列表及其对应的ID：%v，用户会提供一篇企业年报或绿色专利，你需要从关键词列表里面找到年报或专利涉及的所有关键词对应的ID，然后以JSON字符串数组的形式输出。示例输入：本专利报告聚焦绿色技术领域，涵盖清洁能源开发、节能降耗措施及碳中和路径。报告介绍了一项创新型太阳能转化系统，能有效提升能源利用率，显著降低传统能耗，实现低碳运行。该技术有助于推动绿色转型与可持续发展。示例输出：{"keywords":[1,2,3]}`, string(byteKeywords))
			fmt.Println(SystemPrompt)

			completion, err := client.Chat.Completions.New(
				context.Background(),
				openai.ChatCompletionNewParams{
					Model: "deepseek-chat",
					Messages: []openai.ChatCompletionMessageParamUnion{
						openai.SystemMessage(SystemPrompt),
						openai.UserMessage(content),
					},
					ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
						OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
							Type: "json_object",
						},
					},
				},
			)
			if err != nil {
				c.JSON(500, Resp("请求失败", err, nil))
				return
			}

			var data DeepSeekKeywords

			if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &data); err != nil {
				c.JSON(500, bc.Resp("反序列化失败", err, nil))
				return
			}

			var keywords []Keyword
			if err := bc.DB.Where("id IN ?", data.Keywords).Find(&keywords).Error; err != nil {
				c.JSON(500, bc.Resp("查询目标关键词失败", err, nil))
				return
			}

			if err := bc.DB.Model(&patent).Association("Keywords").Replace(&keywords); err != nil {
				c.JSON(500, bc.Resp("替换关键词失败", err, nil))
				return
			}

			c.JSON(200, bc.Resp("请求成功", nil, keywords))
		},
	))

	r.GET("/generate/patent/:id", p.Preload(
		&p.Config[ExtractReq]{
			Base: bc,
			Bind: &p.BindConfig{Param: true},
		},
		func(c *gin.Context, u *User, r *ExtractReq) {

			var patent Patent
			if err := bc.DB.Where("id = ?", r.ID).First(&patent).Error; err != nil {
				c.JSON(500, bc.Resp("获取专利失败", err, nil))
				return
			}

			content, err := GetPDFContent(filepath.Join(cfg.Path.Patent, patent.File))
			if err != nil {
				c.JSON(500, bc.Resp("提取PDF内容失败", err, nil))
				return
			}

			SystemPrompt := `用户会提供一篇企业年报或绿色专利，你需要为专利或年报提出建议，字数100字左右。`

			completion, err := client.Chat.Completions.New(
				context.Background(),
				openai.ChatCompletionNewParams{
					Model: "deepseek-chat",
					Messages: []openai.ChatCompletionMessageParamUnion{
						openai.SystemMessage(SystemPrompt),
						openai.UserMessage(content),
					},
				},
			)
			if err != nil {
				c.JSON(500, Resp("请求失败", err, nil))
				return
			}

			result := completion.Choices[0].Message.Content

			if err := bc.DB.Create(&Suggestion{
				PatentID: &patent.ID,
				UserID:   1,
				SuggestionDTO: &SuggestionDTO{
					Title:   time.Now().Format("2006年01月02日") + strconv.Itoa(rand.Intn(10000)),
					Content: result,
				},
			}).Error; err != nil {
				c.JSON(500, bc.Resp("存储结果到数据库失败", err, nil))
				return
			}

			c.JSON(200, bc.Resp("请求成功", nil, result))
		},
	))

	r.GET("/extract/report/:id", p.Preload(
		&p.Config[ExtractReq]{
			Base: bc,
			Bind: &p.BindConfig{Param: true},
		},
		func(c *gin.Context, u *User, r *ExtractReq) {

			var report Report
			if err := bc.DB.Where("id = ?", r.ID).First(&report).Error; err != nil {
				c.JSON(500, bc.Resp("获取年报失败", err, nil))
				return
			}

			content, err := GetPDFContent(filepath.Join(cfg.Path.Report, report.File))
			if err != nil {
				c.JSON(500, bc.Resp("提取PDF内容失败", err, nil))
				return
			}

			var allKeywords []Keyword
			if err := bc.DB.Select("id", "value").Find(&allKeywords).Error; err != nil {
				c.JSON(500, bc.Resp("查询所有关键词失败", err, nil))
				return
			}

			byteKeywords, err := json.Marshal(&allKeywords)
			if err != nil {
				c.JSON(500, bc.Resp("序列化失败", err, nil))
				return
			}

			fmt.Println(content)

			SystemPrompt := fmt.Sprintf(`这是关键词列表及其对应的ID：%v，用户会提供一篇企业年报或绿色专利，你需要从关键词列表里面找到年报或专利涉及的所有关键词对应的ID，然后以JSON字符串数组的形式输出。示例输入：本专利报告聚焦绿色技术领域，涵盖清洁能源开发、节能降耗措施及碳中和路径。报告介绍了一项创新型太阳能转化系统，能有效提升能源利用率，显著降低传统能耗，实现低碳运行。该技术有助于推动绿色转型与可持续发展。示例输出：{"keywords":[1,2,3]}`, string(byteKeywords))

			completion, err := client.Chat.Completions.New(
				context.Background(),
				openai.ChatCompletionNewParams{
					Model: "deepseek-chat",
					Messages: []openai.ChatCompletionMessageParamUnion{
						openai.SystemMessage(SystemPrompt),
						openai.UserMessage(content),
					},
					ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
						OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
							Type: "json_object",
						},
					},
				},
			)
			if err != nil {
				c.JSON(500, Resp("请求失败", err, nil))
				return
			}

			var data DeepSeekKeywords

			if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &data); err != nil {
				c.JSON(500, bc.Resp("反序列化失败", err, nil))
				return
			}

			var keywords []Keyword
			if err := bc.DB.Where("id IN ?", data.Keywords).Find(&keywords).Error; err != nil {
				c.JSON(500, bc.Resp("查询目标关键词失败", err, nil))
				return
			}

			if err := bc.DB.Model(&report).Association("Keywords").Replace(&keywords); err != nil {
				c.JSON(500, bc.Resp("替换关键词失败", err, nil))
				return
			}

			c.JSON(200, bc.Resp("请求成功", nil, keywords))
		},
	))

	r.GET("/generate/report/:id", p.Preload(
		&p.Config[ExtractReq]{
			Base: bc,
			Bind: &p.BindConfig{Param: true},
		},
		func(c *gin.Context, u *User, r *ExtractReq) {

			var report Report
			if err := bc.DB.Where("id = ?", r.ID).First(&report).Error; err != nil {
				c.JSON(500, bc.Resp("获取年报失败", err, nil))
				return
			}

			content, err := GetPDFContent(filepath.Join(cfg.Path.Report, report.File))
			if err != nil {
				c.JSON(500, bc.Resp("提取PDF内容失败", err, nil))
				return
			}

			SystemPrompt := `用户会提供一篇企业年报或绿色专利，你需要为专利或年报提出建议，字数100字左右。`

			completion, err := client.Chat.Completions.New(
				context.Background(),
				openai.ChatCompletionNewParams{
					Model: "deepseek-chat",
					Messages: []openai.ChatCompletionMessageParamUnion{
						openai.SystemMessage(SystemPrompt),
						openai.UserMessage(content),
					},
				},
			)
			if err != nil {
				c.JSON(500, Resp("请求失败", err, nil))
				return
			}

			result := completion.Choices[0].Message.Content

			if err := bc.DB.Create(&Suggestion{
				ReportID: &report.ID,
				UserID:   1,
				SuggestionDTO: &SuggestionDTO{
					Title:   time.Now().Format("2006年01月02日") + strconv.Itoa(rand.Intn(10000)),
					Content: result,
				},
			}).Error; err != nil {
				c.JSON(500, bc.Resp("存储结果到数据库失败", err, nil))
				return
			}

			c.JSON(200, bc.Resp("请求成功", nil, result))
		},
	))

}

func GetPDFContent(path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	reader, err := pdf.NewReader(file, fileInfo.Size())
	if err != nil {
		return "", err
	}

	text, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(text)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

package parser

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// HTTPResponse represents the structure of an HTTP response.
type HTTPResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	// SourceDouYin represents the source "douyin" (抖音)
	SourceDouYin = "douyin"
	// SourceKuaiShou represents the source "kuaishou" (快手)
	SourceKuaiShou = "kuaishou"
	// SourceTiktok represents the source "tiktok"
	SourceTiktok = "tiktok"
)

// VideoParseInfo 视频解析信息
type VideoParseInfo struct {
	Author       Author          `json:"author"`
	AuthorSource Source          `json:"author_src"`
	Video        SourceVideoInfo `json:"video"`
	VideoSrc     Source          `json:"video_src"`
}

// Source  represents the parsed information based on .
type Source struct {
	Data any `json:"data"`
}

// SourceVideoInfo represents the video information.
type SourceVideoInfo struct {
	Title      string   `json:"title"`     // 描述
	VideoURL   string   `json:"video_url"` // 视频播放地址
	MusicURL   string   `json:"music_url"` // 音乐播放地址
	CoverURL   string   `json:"cover_url"` // 视频封面地址
	Images     []string `json:"images"`    // 图集图片地址列表
	SourceData Source   `json:"video_src"`
}

// Author represents the author information.
type Author struct {
	UID        string      `json:"uid"`         // 作者id
	Name       string      `json:"name"`        // 作者名称
	Avatar     string      `json:"avatar"`      // 作者头像
	SourceData interface{} `json:"source_data"` // 作者
}

// VideoParser is a parser for video information.
type VideoParser struct {
	VideoShareURLParser videoShareURLParser
	VideoIDParser       videoIDParser
}
type videoShareURLParser interface {
	ParseByShareURL(shareURL string) (*VideoParseInfo, error)
}
type videoIDParser interface {
	ParseByShareID(shareID string) (*VideoParseInfo, error)
}

var videoParserMap = map[string]VideoParser{
	SourceDouYin: {
		VideoIDParser:       douYin{},
		VideoShareURLParser: douYin{},
	},
	SourceKuaiShou: {
		VideoIDParser:       douYin{},
		VideoShareURLParser: douYin{},
	},
	SourceTiktok: {
		VideoIDParser:       tiktok{},
		VideoShareURLParser: tiktok{},
	},
}

// Parse parses video information from the given context.
func Parse(c *gin.Context) (*VideoParseInfo, error) {

	query := c.Query("query")
	source := c.Query("source")
	sourceParser, exist := videoParserMap[source]
	if !exist {
		return &VideoParseInfo{}, fmt.Errorf(fmt.Sprintf("Unknown parser %s", source))
	}
	return sourceParser.VideoShareURLParser.ParseByShareURL(query)
	// return ParseVideoID("7483514514890067219")
	// return ParseVideoShareURL("https://www.tiktok.com/@naploes/video/7480531425259834642?is_from_webapp=1&sender_device=pc")
}

package parser

const (
	SourceDouYin   = "douyin"   // 抖音
	SourceKuaiShou = "kuaishou" // 快手

)

// VideoParseInfo 视频解析信息
type VideoParseInfo struct {
	Author struct {
		Uid    string `json:"uid"`    // 作者id
		Name   string `json:"name"`   // 作者名称
		Avatar string `json:"avatar"` // 作者头像
	} `json:"author"`
	Title      string   `json:"title"`     // 描述
	VideoUrl   string   `json:"video_url"` // 视频播放地址
	MusicUrl   string   `json:"music_url"` // 音乐播放地址
	CoverUrl   string   `json:"cover_url"` // 视频封面地址
	Images     []string `json:"images"`    // 图集图片地址列表
	SourceData string   `json:"source_data"`
}

type VideoParser struct {
	VideoShareUrlParser videoShareUrlParser
	VideoIdParser       videoIdParser
}
type videoShareUrlParser interface {
	ParseByShareUrl(shareUrl string) (*VideoParseInfo, error)
}
type videoIdParser interface {
	ParseByShareId(shareId string) (*VideoParseInfo, error)
}

var videoParserMap = map[string]VideoParser{
	SourceDouYin: {
		VideoIdParser:       douYin{},
		VideoShareUrlParser: douYin{},
	},
	SourceKuaiShou: {
		VideoIdParser:       douYin{},
		VideoShareUrlParser: douYin{},
	},
}

func ParseVideoShareUrl(shareUrl string) (*VideoParseInfo, error) {
	source := ""

	parser := videoParserMap[source].VideoShareUrlParser
	return parser.ParseByShareUrl(shareUrl)
}

func ParseVideoId(shareId string) (*VideoParseInfo, error) {
	source := ""
	parser := videoParserMap[source].VideoIdParser
	return parser.ParseByShareId(shareId)
}

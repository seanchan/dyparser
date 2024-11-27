package parser

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type douYin struct{}

func (d douYin) ParseByShareId(videoId string) (*VideoParseInfo, error) {
	reqUrl := fmt.Sprintf("https://www.iesdouyin.com/share/video/%s", videoId)
	client := resty.New()
	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1 Edg/122.0.0.0").
		Get(reqUrl)
	if err != nil {
		return nil, err
	}

	regex, _ := regexp.Compile(`window._ROUTER_DATA\s*=\s*(.*?)</script>`)
	findRes := regex.FindSubmatch(res.Body())
	if len(findRes) < 2 {
		return nil, errors.New("find failed")
	}
	jsonBytes := bytes.TrimSpace(findRes[1])
	data := gjson.GetBytes(jsonBytes, "loaderData.video_(id)/page.videoInfoRes.item_list.0")
	if !data.Exists() {
		filterObj := gjson.GetBytes(jsonBytes, fmt.Sprintf(`loaderData.video_(id)/page.videoInfoRes.filter_list.#(aweme_id=="%s")`, videoId))
		return nil, fmt.Errorf(
			"get video info fail: %s - %s",
			filterObj.Get("filter_reason"),
			filterObj.Get("detail_msg"))
	}
	videoParseInfo := &VideoParseInfo{
		SourceData: data.String(),
	}
	return videoParseInfo, nil

}

func (d douYin) ParseByShareUrl(shareUrl string) (*VideoParseInfo, error) {
	client := resty.New()
}

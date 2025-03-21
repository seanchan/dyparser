package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/go-resty/resty/v2"
)

type douYin struct{}

func (d douYin) ParseByShareID(videoId string) (*VideoParseInfo, error) {
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
	jsonVal := make(map[string]interface{})
	json.Unmarshal(jsonBytes, &jsonVal)

	videoParseInfo := &VideoParseInfo{
		VideoSrc: Source{Data: jsonVal},
	}
	return videoParseInfo, nil

}

func (d douYin) ParseByShareURL(shareUrl string) (*VideoParseInfo, error) {
	client := resty.New()

	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1 Edg/122.0.0.0").
		Get(shareUrl)
	if err != nil {
		return nil, err
	}
	regex, _ := regexp.Compile(`window._ROUTER_DATA\s*=\s*(.*?)</script>`)

	findRes := regex.FindSubmatch(res.Body())
	if len(findRes) < 2 {
		return nil, errors.New("find failed")
	}
	jsonBytes := bytes.TrimSpace(findRes[1])
	jsonVal := make(map[string]interface{})
	json.Unmarshal(jsonBytes, &jsonVal)

	videoParseInfo := &VideoParseInfo{
		VideoSrc: Source{Data: jsonVal},
	}
	return videoParseInfo, nil
}

package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type tiktok struct{}

func (d tiktok) ParseByShareID(videoID string) (*VideoParseInfo, error) {
	reqURL := fmt.Sprintf("https://www.tiktok.com/_/%s", videoID)
	client := resty.New()
	proxyURL := viper.GetString("proxy.url")
	if proxyURL != "" {
		client.SetProxy(proxyURL)
	}
	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1 Edg/122.0.0.0").
		Get(reqURL)
	if err != nil {
		return nil, err
	}

	regex, _ := regexp.Compile(`<script id="__UNIVERSAL_DATA_FOR_REHYDRATION__" type="application/json">\s*(.*?)</script>`)
	findRes := regex.FindSubmatch(res.Body())
	log.Println("found regex", findRes)
	if len(findRes) < 2 {
		return nil, errors.New("find failed")
	}
	jsonBytes := bytes.TrimSpace(findRes[1])
	jsonVal := make(map[string]any)

	json.Unmarshal(jsonBytes, &jsonVal)

	scope := jsonVal["__DEFAULT_SCOPE__"].(map[string]any)
	videoParseInfo := &VideoParseInfo{
		VideoSrc: Source{Data: scope},
	}
	return videoParseInfo, nil

}

func (d tiktok) ParseByShareURL(shareURL string) (*VideoParseInfo, error) {
	client := resty.New()
	proxyURL := viper.GetString("proxy.url")
	if proxyURL != "" {
		client.SetProxy(proxyURL)
	}

	res, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1 Edg/122.0.0.0").
		Get(shareURL)
	if err != nil {
		return nil, err
	}
	regex, _ := regexp.Compile(`<script id="__UNIVERSAL_DATA_FOR_REHYDRATION__" type="application/json">\s*(.*?)</script>`)

	findRes := regex.FindSubmatch(res.Body())
	// log.Println("found regex", findRes)
	if len(findRes) < 2 {
		return nil, errors.New("find failed")
	}
	jsonBytes := bytes.TrimSpace(findRes[1])
	log.Println("found json", string(jsonBytes))
	jsonVal := make(map[string]any)
	json.Unmarshal(jsonBytes, &jsonVal)

	scope := jsonVal["__DEFAULT_SCOPE__"].(map[string]any)
	videoDetail := scope["webapp.reflow.video.detail"].(map[string]any)
	itemInfo := videoDetail["itemInfo"].(map[string]any)
	videoParseInfo := &VideoParseInfo{
		VideoSrc: Source{Data: itemInfo},
	}
	return videoParseInfo, nil
}

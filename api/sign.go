package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/misc/sign/wbi.md

var wbiKeys WbiKeys

func WbiKeysSignString(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	err = wbiKeys.Sign(parsedURL)
	if err != nil {
		return "", err
	}

	return parsedURL.String(), nil
}

// Sign 为链接签名
func WbiKeysSign(u *url.URL) error {
	return wbiKeys.Sign(u)
}

// Update 无视过期时间更新
func WbiKeysUpdate() error {
	return wbiKeys.Update()
}

func WbiKeysGet() (wk WbiKeys, err error) {
	if err = wk.update(false); err != nil {
		return WbiKeys{}, err
	}
	return wbiKeys, nil
}

var mixinKeyEncTab = [...]int{
	46, 47, 18, 2, 53, 8, 23, 32,
	15, 50, 10, 31, 58, 3, 45, 35,
	27, 43, 5, 49, 33, 9, 42, 19,
	29, 28, 14, 39, 12, 38, 41, 13,
	37, 48, 7, 16, 24, 55, 40, 61,
	26, 17, 0, 1, 60, 51, 30, 4,
	22, 25, 54, 21, 56, 59, 6, 63,
	57, 62, 11, 36, 20, 34, 44, 52,
}

func removeUnwantedChars(v url.Values, chars ...byte) url.Values {
	b := []byte(v.Encode())
	for _, c := range chars {
		b = bytes.ReplaceAll(b, []byte{c}, nil)
	}
	s, err := url.ParseQuery(string(b))
	if err != nil {
		panic(err)
	}
	return s
}

type Nav struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		WbiImg struct {
			ImgUrl string `json:"img_url"`
			SubUrl string `json:"sub_url"`
		} `json:"wbi_img"`

		// ......
	} `json:"data"`
}

type WbiKeys struct {
	Img            string
	Sub            string
	Mixin          string
	lastUpdateTime time.Time
}

// Sign 为链接签名
func (wk *WbiKeys) Sign(u *url.URL) (err error) {
	if err = wk.update(false); err != nil {
		return err
	}

	values := u.Query()

	values = removeUnwantedChars(values, '!', '\'', '(', ')', '*') // 必要性存疑?

	values.Set("wts", strconv.FormatInt(time.Now().Unix(), 10))

	// [url.Values.Encode] 内会对参数排序,
	// 且遍历 map 时本身就是无序的
	hash := md5.Sum([]byte(values.Encode() + wk.Mixin)) // Calculate w_rid
	values.Set("w_rid", hex.EncodeToString(hash[:]))
	u.RawQuery = values.Encode()
	return nil
}

// Update 无视过期时间更新
func (wk *WbiKeys) Update() (err error) {
	return wk.update(true)
}

// update 按需更新
func (wk *WbiKeys) update(purge bool) error {
	if !purge && time.Since(wk.lastUpdateTime) < time.Hour {
		return nil
	}

	// 测试下来不用修改 header 也能过
	resp, err := http.Get("https://api.bilibili.com/x/web-interface/nav")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	nav := Nav{}
	err = json.Unmarshal(body, &nav)
	if err != nil {
		return err
	}

	if nav.Code != 0 && nav.Code != -101 { // -101 未登录时也会返回两个 key
		return fmt.Errorf("unexpected code: %d, message: %s", nav.Code, nav.Message)
	}
	img := nav.Data.WbiImg.ImgUrl
	sub := nav.Data.WbiImg.SubUrl
	if img == "" || sub == "" {
		return fmt.Errorf("empty image or sub url: %s", body)
	}

	// https://i0.hdslb.com/bfs/wbi/7cd084941338484aae1ad9425b84077c.png
	imgParts := strings.Split(img, "/")
	subParts := strings.Split(sub, "/")

	// 7cd084941338484aae1ad9425b84077c.png
	imgPng := imgParts[len(imgParts)-1]
	subPng := subParts[len(subParts)-1]

	// 7cd084941338484aae1ad9425b84077c
	wbiKeys.Img = strings.TrimSuffix(imgPng, ".png")
	wbiKeys.Sub = strings.TrimSuffix(subPng, ".png")

	wbiKeys.mixin()
	wbiKeys.lastUpdateTime = time.Now()
	return nil
}

func (wk *WbiKeys) mixin() {
	var mixin [32]byte
	wbi := wk.Img + wk.Sub
	for i := range mixin { // for i := 0; i < len(mixin); i++ {
		mixin[i] = wbi[mixinKeyEncTab[i]]
	}
	wk.Mixin = string(mixin[:])
}

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type DanmakuRequest struct {
	Msg      string
	RoomID   string
	Bubble   string
	Color    string
	FontSize string
	Mode     string
	DmType   string
}

type SendDanmakuResp struct {
	Code int `json:"code"`
	Data struct {
		ModeInfo struct {
			Mode           int    `json:"mode"`
			ShowPlayerType int    `json:"show_player_type"`
			Extra          string `json:"extra"`
		} `json:"mode_info"`
	} `json:"data"`
	Message string `json:"message"`
	Msg     string `json:"msg"`
}

type BiliVerify struct {
	Csrf     string
	SessData string
}

// SendDanmaku https://api.live.bilibili.com/msg/send
func SendDanmaku(d *DanmakuRequest, v *BiliVerify) (*SendDanmakuResp, error) {
	client := &http.Client{}
	result := &SendDanmakuResp{}
	form := url.Values{
		"bubble":     {d.Bubble},
		"color":      {d.Color},
		"fontsize":   {d.FontSize},
		"mode":       {d.Mode},
		"msg":        {d.Msg},
		"roomid":     {d.RoomID},
		"csrf":       {v.Csrf},
		"csrf_token": {v.Csrf},
		"rnd":        {"1"},
	}
	// dm_type 为 1 时，发送的是表情弹幕
	if d.DmType != "" {
		form.Add("dm_type", d.DmType)
	}
	req, err := http.NewRequest("POST", "https://api.live.bilibili.com/msg/send", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", fmt.Sprintf("bili_jct=%s;SESSDATA=%s", v.Csrf, v.SessData))
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}

func SendDefaultDanmaku(roomID string, message string, verify *BiliVerify) (*SendDanmakuResp, error) {
	req := &DanmakuRequest{
		Msg:      message,
		RoomID:   roomID,
		Bubble:   "0",
		Color:    "16777215",
		FontSize: "25",
		Mode:     "1",
		DmType:   "1",
	}
	return SendDanmaku(req, verify)
}

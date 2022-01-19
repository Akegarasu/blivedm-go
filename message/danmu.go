package message

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Danmuku struct {
	Sender  User
	Content string
	Extra   Extra
}

type Extra struct {
	SendFromMe     bool   `json:"send_from_me"`
	Mode           int    `json:"mode"`
	Color          int    `json:"color"`
	DmType         int    `json:"dm_type"`
	FontSize       int    `json:"font_size"`
	PlayerMode     int    `json:"player_mode"`
	ShowPlayerType int    `json:"show_player_type"`
	Content        string `json:"content"`
	UserHash       string `json:"user_hash"`
	EmoticonUnique string `json:"emoticon_unique"`
	Direction      int    `json:"direction"`
	PkDirection    int    `json:"pk_direction"`
	SpaceType      string `json:"space_type"`
	SpaceUrl       string `json:"space_url"`
}

func (d *Danmuku) Parse(data []byte) {
	sb := bytes.NewBuffer(data).String()
	info := gjson.Get(sb, "info")
	d.Content = info.Get("1").String()
	d.Sender = User{
		Uid:   int(info.Get("2.0").Int()),
		Uname: info.Get("2.1").String(),
		Medal: Medal{
			Name:  info.Get("3.1").String(),
			Level: int(info.Get("3.0").Int()),
			Up:    info.Get("3.2").String(),
		},
	}
	ext := new(Extra)
	err := json.Unmarshal([]byte(info.Get("0.15.extra").String()), ext)
	if err != nil {
		log.Error("parse danmuku extra failed")
	}
	d.Extra = *ext
}

package message

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	TextDanmaku = iota
	EmoticonDanmaku
)

type (
	Danmaku struct {
		Sender    *User
		Content   string
		Extra     *Extra
		Emoticon  *Emoticon
		Type      int
		Timestamp int64
	}

	Extra struct {
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
	Emoticon struct {
		BulgeDisplay   int    `json:"bulge_display"`
		EmoticonUnique string `json:"emoticon_unique"`
		Height         int    `json:"height"`
		InPlayerArea   int    `json:"in_player_area"`
		IsDynamic      int    `json:"is_dynamic"`
		Url            string `json:"url"`
		Width          int    `json:"width"`
	}
)

func (d *Danmaku) Parse(data []byte) {
	sb := bytes.NewBuffer(data).String()
	info := gjson.Get(sb, "info")
	ext := new(Extra)
	emo := new(Emoticon)
	err := json.Unmarshal([]byte(info.Get("0.15.extra").String()), ext)
	if err != nil {
		log.Error("parse danmuku extra failed")
	}
	err = json.Unmarshal([]byte(info.Get("0.13").String()), emo)
	if err != nil {
		log.Error("parse danmuku emoticon failed")
	}
	d.Content = info.Get("1").String()
	d.Sender = &User{
		Uid:   int(info.Get("2.0").Int()),
		Uname: info.Get("2.1").String(),
		Medal: &Medal{
			Name:  info.Get("3.1").String(),
			Level: int(info.Get("3.0").Int()),
			Up:    info.Get("3.2").String(),
		},
	}
	d.Extra = ext
	d.Emoticon = emo
	d.Type = int(info.Get("0.12").Int())
	d.Timestamp = info.Get("0.4").Int()
}

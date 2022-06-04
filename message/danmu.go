package message

import (
	"github.com/Akegarasu/blivedm-go/utils"
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
		Raw       string
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
	CommonNoticeDanmaku struct {
		ContentSegments []struct {
			FontColor string `json:"font_color"`
			Text      string `json:"text"`
			Type      int    `json:"type"`
		} `json:"content_segments"`
		Dmscore   int   `json:"dmscore"`
		Terminals []int `json:"terminals"`
	}
)

func (d *Danmaku) Parse(data []byte) {
	sb := utils.BytesToString(data)
	info := gjson.Parse(sb).Get("info")
	ext := new(Extra)
	emo := new(Emoticon)
	err := utils.UnmarshalString(info.Get("0.15.extra").String(), ext)
	if err != nil {
		log.Error("parse danmuku extra failed")
	}
	err = utils.UnmarshalString(info.Get("0.13").String(), emo)
	if err != nil {
		log.Error("parse danmuku emoticon failed")
	}
	d.Content = info.Get("1").String()
	d.Sender = &User{
		Uid:          int(info.Get("2.0").Int()),
		Uname:        info.Get("2.1").String(),
		Admin:        info.Get("2.2").Bool(),
		Urank:        int(info.Get("2.5").Int()),
		MobileVerify: info.Get("2.6").Bool(),
		GuardLevel:   int(info.Get("7").Int()),
		Medal: &Medal{
			Level:    int(info.Get("3.0").Int()),
			Name:     info.Get("3.1").String(),
			UpName:   info.Get("3.2").String(),
			UpRoomId: int(info.Get("3.3").Int()),
			Color:    int(info.Get("3.4").Int()),
			UpUid:    int(info.Get("3.12").Int()),
		},
	}
	d.Extra = ext
	d.Emoticon = emo
	d.Type = int(info.Get("0.12").Int())
	d.Timestamp = info.Get("0.4").Int()
	d.Raw = sb
}

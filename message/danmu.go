package message

import (
	"github.com/Akegarasu/blivedm-go/pb"
	"github.com/Akegarasu/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/proto"
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
	parsed := gjson.Parse(sb)
	info := parsed.Get("info")

	ext := new(Extra)
	emo := new(Emoticon)
	err := utils.UnmarshalStr(info.Get("0.15.extra").String(), ext)
	if err != nil {
		log.Error("parse danmaku extra failed")
	}
	err = utils.UnmarshalStr(info.Get("0.13").String(), emo)
	if err != nil {
		log.Error("parse danmaku emoticon failed")
	}
	i2 := info.Get("2")
	i3 := info.Get("3")
	d.Content = info.Get("1").String()
	d.Sender = &User{
		Uid:          int(i2.Get("0").Int()),
		Uname:        i2.Get("1").String(),
		Admin:        i2.Get("2").Bool(),
		Urank:        int(i2.Get("5").Int()),
		MobileVerify: i2.Get("6").Bool(),
		GuardLevel:   int(info.Get("7").Int()),
		Medal: &Medal{
			Level:    int(i3.Get("0").Int()),
			Name:     i3.Get("1").String(),
			UpName:   i3.Get("2").String(),
			UpRoomId: int(i3.Get("3").Int()),
			Color:    int(i3.Get("4").Int()),
			UpUid:    int(i3.Get("12").Int()),
		},
	}
	d.Extra = ext
	d.Emoticon = emo
	d.Type = int(info.Get("0.12").Int())
	d.Timestamp = info.Get("0.4").Int()
	d.Raw = sb

	dmv2Content := parsed.Get("dm_v2").String()
	if dmv2Content != "" {
		decoded, _ := utils.B64Decode(dmv2Content)
		dmv2 := new(pb.Dm)

		err := proto.Unmarshal(decoded, dmv2)
		if err != nil {
			return
		}
		d.Content = dmv2.Content
	}

}

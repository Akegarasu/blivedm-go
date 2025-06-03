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
		SendFromMe            bool   `json:"send_from_me"`
		MasterPlayerHidden    bool   `json:"master_player_hidden"`
		Mode                  int    `json:"mode"`
		Color                 int    `json:"color"`
		DmType                int    `json:"dm_type"`
		FontSize              int    `json:"font_size"`
		PlayerMode            int    `json:"player_mode"`
		ShowPlayerType        int    `json:"show_player_type"`
		Content               string `json:"content"`
		UserHash              string `json:"user_hash"`
		EmoticonUnique        string `json:"emoticon_unique"`
		BulgeDisplay          int    `json:"bulge_display"`
		RecommendScore        int    `json:"recommend_score"`
		MainStateDmColor      string `json:"main_state_dm_color"`
		ObjectiveStateDmColor string `json:"objective_state_dm_color"`
		Direction             int    `json:"direction"`
		PkDirection           int    `json:"pk_direction"`
		QuartetDirection      int    `json:"quartet_direction"`
		AnniversaryCrowd      int    `json:"anniversary_crowd"`
		YeahSpaceType         string `json:"yeah_space_type"`
		YeahSpaceURL          string `json:"yeah_space_url"`
		JumpToURL             string `json:"jump_to_url"`
		SpaceType             string `json:"space_type"`
		SpaceURL              string `json:"space_url"`
		// Animation             any    `json:"animation"`
		// Emots                 any    `json:"emots"`
		IsAudited bool   `json:"is_audited"`
		IDStr     string `json:"id_str"`
		// Icon                  any    `json:"icon"`
		ShowReply       bool   `json:"show_reply"`
		ReplyMid        int    `json:"reply_mid"`
		ReplyUname      string `json:"reply_uname"`
		ReplyUnameColor string `json:"reply_uname_color"`
		ReplyIsMystery  bool   `json:"reply_is_mystery"`
		ReplyTypeEnum   int    `json:"reply_type_enum"`
		HitCombo        int    `json:"hit_combo"`
		EsportsJumpURL  string `json:"esports_jump_url"`
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
	//扩展字段
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
	d.Content = info.Get("1").String() //弹幕内容

	d.Sender = &User{ //用户信息
		Uid:          int(i2.Get("0").Int()), //用户uid
		Uname:        i2.Get("1").String(),   //用户昵称
		UserLevel:    i2.Get("16.0").Int(),   //用户等级
		Admin:        i2.Get("2").Bool(),     //是否为管理者
		Urank:        int(i2.Get("5").Int()),
		MobileVerify: i2.Get("6").Bool(),       //是否绑定手机
		GuardLevel:   int(info.Get("7").Int()), //舰队等级
		//勋章信息
		Medal: &Medal{
			Level:    int(i3.Get("0").Int()),  //勋章等级
			Name:     i3.Get("1").String(),    //勋章名称
			UpName:   i3.Get("2").String(),    //勋章上主播昵称
			UpRoomId: int(i3.Get("3").Int()),  //上主播房间id
			Color:    int(i3.Get("4").Int()),  //勋章颜色
			UpUid:    int(i3.Get("12").Int()), //上主播uid
		},
	}
	d.Extra = ext
	d.Emoticon = emo
	d.Type = int(info.Get("0.12").Int()) //弹幕类型
	d.Timestamp = info.Get("0.4").Int()  //时间戳
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

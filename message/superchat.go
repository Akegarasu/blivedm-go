package message

import (
	"github.com/Akegarasu/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// SuperChat
// message_jpn: 消息日文翻译（目前只出现在SUPER_CHAT_MESSAGE_JPN）
// id_: str，消息ID，删除时用
type SuperChat struct {
	BackgroundBottomColor string  `json:"background_bottom_color"` //底部背景色
	BackgroundColor       string  `json:"background_color"`        //背景色
	BackgroundColorEnd    string  `json:"background_color_end"`
	BackgroundColorStart  string  `json:"background_color_start"`
	BackgroundIcon        string  `json:"background_icon"`        //背景图标
	BackgroundImage       string  `json:"background_image"`       //背景图
	BackgroundPriceColor  string  `json:"background_price_color"` //背景价格颜色
	ColorPoint            float64 `json:"color_point"`
	Dmscore               int     `json:"dmscore"`
	EndTime               int     `json:"end_time"` //结束时间戳
	Gift                  struct {
		GiftId   int    `json:"gift_id"`   //礼物ID
		GiftName string `json:"gift_name"` //礼物名
		Num      int    `json:"num"`
	} `json:"gift"`
	Id          int `json:"id"`
	IsRanked    int `json:"is_ranked"`
	IsSendAudit int `json:"is_send_audit"`
	MedalInfo   struct {
		AnchorRoomid     int    `json:"anchor_roomid"`
		AnchorUname      string `json:"anchor_uname"`
		GuardLevel       int    `json:"guard_level"` // 舰队等级，0:非舰队，1:总督，2:提督，3:舰长
		IconId           int    `json:"icon_id"`
		IsLighted        int    `json:"is_lighted"`
		MedalColor       string `json:"medal_color"`
		MedalColorBorder int    `json:"medal_color_border"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		Special          string `json:"special"`
		TargetId         int    `json:"target_id"`
	} `json:"medal_info"`
	Message          string `json:"message"` // 消息
	MessageFontColor string `json:"message_font_color"`
	MessageTrans     string `json:"message_trans"`
	Price            int    `json:"price"` // 价格（人民币）
	Rate             int    `json:"rate"`
	StartTime        int    `json:"start_time"` // 开始时间戳
	Time             int    `json:"time"`       //剩余时间
	Token            string `json:"token"`
	TransMark        int    `json:"trans_mark"`
	Ts               int    `json:"ts"`
	Uid              int    `json:"uid"` //用户ID
	UserInfo         struct {
		Face       string `json:"face"` //用户头像URL
		FaceFrame  string `json:"face_frame"`
		GuardLevel int    `json:"guard_level"`
		IsMainVip  int    `json:"is_main_vip"`
		IsSvip     int    `json:"is_svip"`
		IsVip      int    `json:"is_vip"`
		LevelColor string `json:"level_color"`
		Manager    int    `json:"manager"`
		NameColor  string `json:"name_color"`
		Title      string `json:"title"`
		Uname      string `json:"uname"`      //用户名
		UserLevel  int    `json:"user_level"` //用户等级
	} `json:"user_info"`
}

func (s *SuperChat) Parse(data []byte) {
	sb := utils.BytesToString(data)
	sd := gjson.Get(sb, "data").String()
	err := utils.UnmarshalStr(sd, s)
	if err != nil {
		log.Error("parse superchat failed")
	}
}

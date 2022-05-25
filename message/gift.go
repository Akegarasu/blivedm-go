package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Gift struct {
	Action            string      `json:"action"`
	BatchComboId      string      `json:"batch_combo_id"`
	BatchComboSend    interface{} `json:"batch_combo_send"`
	BeatId            string      `json:"beatId"`
	BizSource         string      `json:"biz_source"`
	BlindGift         interface{} `json:"blind_gift"`
	BroadcastId       int         `json:"broadcast_id"`
	CoinType          string      `json:"coin_type"`
	ComboResourcesId  int         `json:"combo_resources_id"`
	ComboSend         interface{} `json:"combo_send"`
	ComboStayTime     int         `json:"combo_stay_time"`
	ComboTotalCoin    int         `json:"combo_total_coin"`
	CritProb          int         `json:"crit_prob"`
	Demarcation       int         `json:"demarcation"`
	DiscountPrice     int         `json:"discount_price"`
	Dmscore           int         `json:"dmscore"`
	Draw              int         `json:"draw"`
	Effect            int         `json:"effect"`
	EffectBlock       int         `json:"effect_block"`
	Face              string      `json:"face"`
	FloatScResourceId int         `json:"float_sc_resource_id"`
	GiftId            int         `json:"giftId"`
	GiftName          string      `json:"giftName"`
	GiftType          int         `json:"giftType"`
	Gold              int         `json:"gold"`
	GuardLevel        int         `json:"guard_level"`
	IsFirst           bool        `json:"is_first"`
	IsSpecialBatch    int         `json:"is_special_batch"`
	Magnification     float64     `json:"magnification"`
	MedalInfo         struct {
		AnchorRoomid     int    `json:"anchor_roomid"`
		AnchorUname      string `json:"anchor_uname"`
		GuardLevel       int    `json:"guard_level"`
		IconId           int    `json:"icon_id"`
		IsLighted        int    `json:"is_lighted"`
		MedalColor       int    `json:"medal_color"`
		MedalColorBorder int    `json:"medal_color_border"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		Special          string `json:"special"`
		TargetId         int    `json:"target_id"`
	} `json:"medal_info"`
	NameColor         string      `json:"name_color"`
	Num               int         `json:"num"`
	OriginalGiftName  string      `json:"original_gift_name"`
	Price             int         `json:"price"`
	Rcost             int         `json:"rcost"`
	Remain            int         `json:"remain"`
	Rnd               string      `json:"rnd"`
	SendMaster        interface{} `json:"send_master"`
	Silver            int         `json:"silver"`
	Super             int         `json:"super"`
	SuperBatchGiftNum int         `json:"super_batch_gift_num"`
	SuperGiftNum      int         `json:"super_gift_num"`
	SvgaBlock         int         `json:"svga_block"`
	TagImage          string      `json:"tag_image"`
	Tid               string      `json:"tid"`
	Timestamp         int         `json:"timestamp"`
	TopList           interface{} `json:"top_list"`
	TotalCoin         int         `json:"total_coin"`
	Uid               int         `json:"uid"`
	Uname             string      `json:"uname"`
}

type ComboSend struct {
	Action         string `json:"action"`
	BatchComboId   string `json:"batch_combo_id"`
	BatchComboNum  int    `json:"batch_combo_num"`
	ComboId        string `json:"combo_id"`
	ComboNum       int    `json:"combo_num"`
	ComboTotalCoin int    `json:"combo_total_coin"`
	Dmscore        int    `json:"dmscore"`
	GiftId         int    `json:"gift_id"`
	GiftName       string `json:"gift_name"`
	GiftNum        int    `json:"gift_num"`
	IsShow         int    `json:"is_show"`
	MedalInfo      struct {
		AnchorRoomid     int    `json:"anchor_roomid"`
		AnchorUname      string `json:"anchor_uname"`
		GuardLevel       int    `json:"guard_level"`
		IconId           int    `json:"icon_id"`
		IsLighted        int    `json:"is_lighted"`
		MedalColor       int    `json:"medal_color"`
		MedalColorBorder int    `json:"medal_color_border"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		Special          string `json:"special"`
		TargetId         int    `json:"target_id"`
	} `json:"medal_info"`
	NameColor  string      `json:"name_color"`
	RUname     string      `json:"r_uname"`
	Ruid       int         `json:"ruid"`
	SendMaster interface{} `json:"send_master"`
	TotalNum   int         `json:"total_num"`
	Uid        int         `json:"uid"`
	Uname      string      `json:"uname"`
}

func (g *Gift) Parse(data []byte) {
	// len("{"cmd":"","data":") == 17 , len('SEND_GIFT') = 9
	d := data[17+9 : len(data)-1]
	err := json.Unmarshal(d, g)
	if err != nil {
		log.Error("parse Gift failed")
	}
}

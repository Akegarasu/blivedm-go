package message

// WidgetBanner
// TODO: widget_list的code不定
type WidgetBanner struct {
	Timestamp  int `json:"timestamp"`
	WidgetList struct {
		Field1 struct {
			Id             int      `json:"id"`
			Title          string   `json:"title"`
			Cover          string   `json:"cover"`
			WebCover       string   `json:"web_cover"`
			TipText        string   `json:"tip_text"`
			TipTextColor   string   `json:"tip_text_color"`
			TipBottomColor string   `json:"tip_bottom_color"`
			JumpUrl        string   `json:"jump_url"`
			Url            string   `json:"url"`
			StayTime       int      `json:"stay_time"`
			Site           int      `json:"site"`
			PlatformIn     []string `json:"platform_in"`
			Type           int      `json:"type"`
			BandId         int      `json:"band_id"`
			SubKey         string   `json:"sub_key"`
			SubData        string   `json:"sub_data"`
			IsAdd          bool     `json:"is_add"`
		} `json:"58"`
	} `json:"widget_list"`
}

type HotRankChanged struct {
	Rank        int    `json:"rank"`
	Trend       int    `json:"trend"`
	Countdown   int    `json:"countdown"`
	Timestamp   int    `json:"timestamp"`
	WebUrl      string `json:"web_url"`
	LiveUrl     string `json:"live_url"`
	BlinkUrl    string `json:"blink_url"`
	LiveLinkUrl string `json:"live_link_url"`
	PcLinkUrl   string `json:"pc_link_url"`
	Icon        string `json:"icon"`
	AreaName    string `json:"area_name"`
	RankDesc    string `json:"rank_desc"`
}

type HotRankChangedV2 HotRankChanged

type HotRankSettlement struct {
	AreaName  string `json:"area_name"`
	CacheKey  string `json:"cache_key"`
	DmMsg     string `json:"dm_msg"`
	Dmscore   int    `json:"dmscore"`
	Face      string `json:"face"`
	Icon      string `json:"icon"`
	Rank      int    `json:"rank"`
	Timestamp int    `json:"timestamp"`
	Uname     string `json:"uname"`
	Url       string `json:"url"`
}

type HotRankSettlementV2 struct {
	Rank      int    `json:"rank"`
	Uname     string `json:"uname"`
	Face      string `json:"face"`
	Timestamp int    `json:"timestamp"`
	Icon      string `json:"icon"`
	AreaName  string `json:"area_name"`
	Url       string `json:"url"`
	CacheKey  string `json:"cache_key"`
	DmMsg     string `json:"dm_msg"`
}

type InteractWord struct {
	Contribution struct {
		Grade int `json:"grade"`
	} `json:"contribution"`
	Dmscore   int `json:"dmscore"`
	FansMedal struct {
		AnchorRoomid     int    `json:"anchor_roomid"`
		GuardLevel       int    `json:"guard_level"`
		IconId           int    `json:"icon_id"`
		IsLighted        int    `json:"is_lighted"`
		MedalColor       int    `json:"medal_color"`
		MedalColorBorder int    `json:"medal_color_border"`
		MedalColorEnd    int    `json:"medal_color_end"`
		MedalColorStart  int    `json:"medal_color_start"`
		MedalLevel       int    `json:"medal_level"`
		MedalName        string `json:"medal_name"`
		Score            int    `json:"score"`
		Special          string `json:"special"`
		TargetId         int    `json:"target_id"`
	} `json:"fans_medal"`
	Identities  []int  `json:"identities"`
	IsSpread    int    `json:"is_spread"`
	MsgType     int    `json:"msg_type"`
	Roomid      int    `json:"roomid"`
	Score       int64  `json:"score"`
	SpreadDesc  string `json:"spread_desc"`
	SpreadInfo  string `json:"spread_info"`
	TailIcon    int    `json:"tail_icon"`
	Timestamp   int    `json:"timestamp"`
	TriggerTime int64  `json:"trigger_time"`
	Uid         int    `json:"uid"`
	Uname       string `json:"uname"`
	UnameColor  string `json:"uname_color"`
}

type OnlineRankCount struct {
	Count int `json:"count"`
}

type LiveInteractiveGame struct {
	Type           int         `json:"type"`
	Uid            int         `json:"uid"`
	Uname          string      `json:"uname"`
	Uface          string      `json:"uface"`
	GiftId         int         `json:"gift_id"`
	GiftName       string      `json:"gift_name"`
	GiftNum        int         `json:"gift_num"`
	Price          int         `json:"price"`
	Paid           bool        `json:"paid"`
	Msg            string      `json:"msg"`
	FansMedalLevel int         `json:"fans_medal_level"`
	GuardLevel     int         `json:"guard_level"`
	Timestamp      int         `json:"timestamp"`
	AnchorLottery  interface{} `json:"anchor_lottery"`
	PkInfo         interface{} `json:"pk_info"`
	AnchorInfo     struct {
		Uid   int    `json:"uid"`
		Uname string `json:"uname"`
		Uface string `json:"uface"`
	} `json:"anchor_info"`
}

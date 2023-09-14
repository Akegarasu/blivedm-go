package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// RoomInfo
// api https://api.live.bilibili.com/room/v1/Room/room_init?id={} response
type RoomInfo struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		RoomId      int   `json:"room_id"`
		ShortId     int   `json:"short_id"`
		Uid         int   `json:"uid"`
		NeedP2P     int   `json:"need_p2p"`
		IsHidden    bool  `json:"is_hidden"`
		IsLocked    bool  `json:"is_locked"`
		IsPortrait  bool  `json:"is_portrait"`
		LiveStatus  int   `json:"live_status"`
		HiddenTill  int   `json:"hidden_till"`
		LockTill    int   `json:"lock_till"`
		Encrypted   bool  `json:"encrypted"`
		PwdVerified bool  `json:"pwd_verified"`
		LiveTime    int64 `json:"live_time"`
		RoomShield  int   `json:"room_shield"`
		IsSp        int   `json:"is_sp"`
		SpecialType int   `json:"special_type"`
	} `json:"data"`
}

// DanmuInfo
// api https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id={}&type=0 response
type DanmuInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Group            string  `json:"group"`
		BusinessId       int     `json:"business_id"`
		RefreshRowFactor float64 `json:"refresh_row_factor"`
		RefreshRate      int     `json:"refresh_rate"`
		MaxDelay         int     `json:"max_delay"`
		Token            string  `json:"token"`
		HostList         []struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			WssPort int    `json:"wss_port"`
			WsPort  int    `json:"ws_port"`
		} `json:"host_list"`
	} `json:"data"`
}

// UserInfo
// api https://api.bilibili.com/x/web-interface/nav
type UserInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		IsLogin       bool   `json:"isLogin"`
		EmailVerified int    `json:"email_verified"`
		Face          string `json:"face"`
		FaceNft       int    `json:"face_nft"`
		FaceNftType   int    `json:"face_nft_type"`
		LevelInfo     struct {
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
			NextExp      int `json:"next_exp"`
		} `json:"level_info"`
		Mid            int     `json:"mid"`
		MobileVerified int     `json:"mobile_verified"`
		Money          float64 `json:"money"`
		Moral          int     `json:"moral"`
		Official       struct {
			Role  int    `json:"role"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
			Type  int    `json:"type"`
		} `json:"official"`
		OfficialVerify struct {
			Type int    `json:"type"`
			Desc string `json:"desc"`
		} `json:"officialVerify"`
		Pendant struct {
			Pid               int    `json:"pid"`
			Name              string `json:"name"`
			Image             string `json:"image"`
			Expire            int    `json:"expire"`
			ImageEnhance      string `json:"image_enhance"`
			ImageEnhanceFrame string `json:"image_enhance_frame"`
		} `json:"pendant"`
		Scores       int    `json:"scores"`
		Uname        string `json:"uname"`
		VipDueDate   int64  `json:"vipDueDate"`
		VipStatus    int    `json:"vipStatus"`
		VipType      int    `json:"vipType"`
		VipPayType   int    `json:"vip_pay_type"`
		VipThemeType int    `json:"vip_theme_type"`
		VipLabel     struct {
			Path                  string `json:"path"`
			Text                  string `json:"text"`
			LabelTheme            string `json:"label_theme"`
			TextColor             string `json:"text_color"`
			BgStyle               int    `json:"bg_style"`
			BgColor               string `json:"bg_color"`
			BorderColor           string `json:"border_color"`
			UseImgLabel           bool   `json:"use_img_label"`
			ImgLabelURIHans       string `json:"img_label_uri_hans"`
			ImgLabelURIHant       string `json:"img_label_uri_hant"`
			ImgLabelURIHansStatic string `json:"img_label_uri_hans_static"`
			ImgLabelURIHantStatic string `json:"img_label_uri_hant_static"`
		} `json:"vip_label"`
		VipAvatarSubscript int    `json:"vip_avatar_subscript"`
		VipNicknameColor   string `json:"vip_nickname_color"`
		Vip                struct {
			Type       int   `json:"type"`
			Status     int   `json:"status"`
			DueDate    int64 `json:"due_date"`
			VipPayType int   `json:"vip_pay_type"`
			ThemeType  int   `json:"theme_type"`
			Label      struct {
				Path                  string `json:"path"`
				Text                  string `json:"text"`
				LabelTheme            string `json:"label_theme"`
				TextColor             string `json:"text_color"`
				BgStyle               int    `json:"bg_style"`
				BgColor               string `json:"bg_color"`
				BorderColor           string `json:"border_color"`
				UseImgLabel           bool   `json:"use_img_label"`
				ImgLabelURIHans       string `json:"img_label_uri_hans"`
				ImgLabelURIHant       string `json:"img_label_uri_hant"`
				ImgLabelURIHansStatic string `json:"img_label_uri_hans_static"`
				ImgLabelURIHantStatic string `json:"img_label_uri_hant_static"`
			} `json:"label"`
			AvatarSubscript    int    `json:"avatar_subscript"`
			NicknameColor      string `json:"nickname_color"`
			Role               int    `json:"role"`
			AvatarSubscriptURL string `json:"avatar_subscript_url"`
			TvVipStatus        int    `json:"tv_vip_status"`
			TvVipPayType       int    `json:"tv_vip_pay_type"`
			TvDueDate          int    `json:"tv_due_date"`
		} `json:"vip"`
		Wallet struct {
			Mid           int `json:"mid"`
			BcoinBalance  int `json:"bcoin_balance"`
			CouponBalance int `json:"coupon_balance"`
			CouponDueTime int `json:"coupon_due_time"`
		} `json:"wallet"`
		HasShop        bool   `json:"has_shop"`
		ShopURL        string `json:"shop_url"`
		AllowanceCount int    `json:"allowance_count"`
		AnswerStatus   int    `json:"answer_status"`
		IsSeniorMember int    `json:"is_senior_member"`
		WbiImg         struct {
			ImgURL string `json:"img_url"`
			SubURL string `json:"sub_url"`
		} `json:"wbi_img"`
		IsJury bool `json:"is_jury"`
	} `json:"data"`
}

func GetUid(cookie string) (int, error) {
	result := &UserInfo{}
	headers := &http.Header{}
	headers.Set("cookie", cookie)
	err := GetJsonWithHeader("https://api.bilibili.com/x/web-interface/nav", headers, result)
	if err != nil {
		return 0, err
	}
	if result.Code != 0 || !result.Data.IsLogin {
		return 0, errors.New(result.Message)
	}
	return result.Data.Mid, nil
}

func GetDanmuInfo(roomID int, cookie string) (*DanmuInfo, error) {
	result := &DanmuInfo{}
	headers := &http.Header{}
	headers.Set("cookie", cookie)
	err := GetJsonWithHeader(fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%d&type=0", roomID), headers, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetRoomInfo(roomID int) (*RoomInfo, error) {
	result := &RoomInfo{}
	err := GetJson(fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/room_init?id=%d", roomID), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetRoomRealID(roomID int) (string, error) {
	res, err := GetRoomInfo(roomID)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(res.Data.RoomId), nil
}

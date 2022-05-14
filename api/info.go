package api

import (
	"fmt"
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

func GetDanmuInfo(roomID string) (*DanmuInfo, error) {
	result := &DanmuInfo{}
	err := GetJson(fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%s&type=0", roomID), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetRoomInfo(roomID string) (*RoomInfo, error) {
	result := &RoomInfo{}
	err := GetJson(fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/room_init?id=%s", roomID), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetRoomRealID(roomID string) (string, error) {
	res, err := GetRoomInfo(roomID)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(res.Data.RoomId), nil
}

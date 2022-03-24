package api

import (
	"encoding/json"
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

func getRoomInfo(roomID string) (*RoomInfo, error) {
	url := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/room_init?id=%s", roomID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result := &RoomInfo{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetRoomRealID(roomID string) (string, error) {
	res, err := getRoomInfo(roomID)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(res.Data.RoomId), nil
}

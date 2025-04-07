package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

// StopLiveRoomList 停止直播房间列表结构体
type StopLiveRoomList struct {
	RoomIdList []int `json:"room_id_list"` // 房间ID列表
}

// Live 直播消息结构体
type LiveStart struct {
	Cmd             string `json:"cmd"`              // 命令
	LiveKey         string `json:"live_key"`         // 直播密钥
	VoiceBackground string `json:"voice_background"` // 语音背景
	SubSessionKey   string `json:"sub_session_key"`  // 子会话密钥
	LivePlatform    string `json:"live_platform"`    // 直播平台
	LiveModel       int    `json:"live_model"`       // 直播模式
	LiveTime        int    `json:"live_time"`        // 直播时间
	Roomid          int    `json:"roomid"`           // 房间ID
}

// 停止直播
type LiveStop struct {
	Cmd      string `json:"cmd"`
	MsgId    string `json:"msg_id"`
	PIsAck   bool   `json:"p_is_ack"`
	PMsgType int    `json:"p_msg_type"`
	Roomid   string `json:"roomid"`
	Round    int    `json:"round"` //开启轮播时存在,轮播状态: 1正在轮播 0未轮播
	SendTime int64  `json:"send_time"`
}

// Preparing 直播准备中消息结构体
type Preparing struct {
	Cmd    string `json:"cmd"`    // 命令
	Roomid string `json:"roomid"` // 房间ID
}

func (l *LiveStart) Parse(data []byte) {
	err := json.Unmarshal(data, l)
	if err != nil {
		log.Error("parse live failed")
	}
}
func (l *LiveStop) Parse(data []byte) {
	err := json.Unmarshal(data, l)
	if err != nil {
		log.Error("parse live failed")
	}
}

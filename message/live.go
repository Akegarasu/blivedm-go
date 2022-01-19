package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type StopLiveRoomList struct {
	RoomIdList []int `json:"room_id_list"`
}

type Live struct {
	Cmd             string `json:"cmd"`
	LiveKey         string `json:"live_key"`
	VoiceBackground string `json:"voice_background"`
	SubSessionKey   string `json:"sub_session_key"`
	LivePlatform    string `json:"live_platform"`
	LiveModel       int    `json:"live_model"`
	LiveTime        int    `json:"live_time"`
	Roomid          int    `json:"roomid"`
}

type Preparing struct {
	Cmd    string `json:"cmd"`
	Roomid string `json:"roomid"`
}

func (l *Live) Parse(data []byte) {
	err := json.Unmarshal(data, l)
	if err != nil {
		log.Error("parse live failed")
	}
}

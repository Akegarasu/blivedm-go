package packet

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type Enter struct {
	UID       int    `json:"uid"`
	RoomID    int    `json:"roomid"`
	ProtoVer  int    `json:"protover"`
	Platform  string `json:"platform"`
	ClientVer string `json:"clientver"`
	//Type      int    `json:"type"`
	//Key       string `json:"key"`
}

// NewEnterPacket 构造进入房间的包
// uid可以为0，key不需要
func NewEnterPacket(uid int, roomID int) []byte {
	ent := &Enter{
		UID:       uid,
		RoomID:    roomID,
		ProtoVer:  1,
		Platform:  "web",
		ClientVer: "1.6.3",
		//Type:      2,
		//Key:       key,
	}
	pkt := NewPlainPacket(RoomEnter, ent.Json())
	upkt := pkt.Build()
	return upkt
}

func (e *Enter) Json() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		log.Error("NewEnterPacket JsonMarshal failed", err)
	}
	return marshal
}

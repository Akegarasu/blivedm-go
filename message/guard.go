package message

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type GuardBuy struct {
	Uid        int    `json:"uid"`
	Username   string `json:"username"`
	GuardLevel int    `json:"guard_level"`
	Num        int    `json:"num"`
	Price      int    `json:"price"`
	GiftId     int    `json:"gift_id"`
	GiftName   string `json:"gift_name"`
	StartTime  int    `json:"start_time"`
	EndTime    int    `json:"end_time"`
}

func (g *GuardBuy) Parse(data []byte) {
	sb := bytes.NewBuffer(data).String()
	sd := gjson.Get(sb, "data").String()
	err := json.Unmarshal([]byte(sd), g)
	if err != nil {
		log.Error("parse GuardBuy failed")
	}
}

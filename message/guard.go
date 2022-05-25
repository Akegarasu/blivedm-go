package message

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
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
	// len("{"cmd":"","data":") = 17 , len('GUARD_BUY') = 9
	d := data[17+9 : len(data)-1]
	err := json.Unmarshal(d, g)
	if err != nil {
		log.Error("parse GuardBuy failed")
	}
}

package message

import (
	"github.com/Akegarasu/blivedm-go/utils"
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
	sb := utils.BytesToString(data)
	sd := gjson.Get(sb, "data").String()
	err := utils.UnmarshalStr(sd, g)
	if err != nil {
		log.Error("parse GuardBuy failed")
	}
}

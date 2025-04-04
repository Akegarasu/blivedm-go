package message

import (
	"github.com/Akegarasu/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// GuardBuy 购买守护消息结构体
type GuardBuy struct {
	Uid        int    `json:"uid"`         // 用户ID
	Username   string `json:"username"`    // 用户名
	GuardLevel int    `json:"guard_level"` // 守护等级
	Num        int    `json:"num"`         // 数量
	Price      int    `json:"price"`       // 价格
	GiftId     int    `json:"gift_id"`     // 礼物ID
	GiftName   string `json:"gift_name"`   // 礼物名称
	StartTime  int    `json:"start_time"`  // 开始时间戳
	EndTime    int    `json:"end_time"`    // 结束时间戳
}

func (g *GuardBuy) Parse(data []byte) {
	sb := utils.BytesToString(data)
	sd := gjson.Get(sb, "data").String()
	err := utils.UnmarshalStr(sd, g)
	if err != nil {
		log.Error("parse GuardBuy failed")
	}
}

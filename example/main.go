package main

import (
	"fmt"

	"github.com/Akegarasu/blivedm-go/api"
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	_ "github.com/Akegarasu/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func main() {
	log.SetLevel(log.DebugLevel)
	c := client.NewClient("23555200", "194484313")
	c.UseDefaultHost()
	//弹幕事件
	c.OnDanmaku(func(danmaku *message.Danmaku) {
		if danmaku.Type == message.EmoticonDanmaku {
			fmt.Printf("[弹幕表情] %s：表情URL： %s\n", danmaku.Sender.Uname, danmaku.Emoticon.Url)
		} else {
			fmt.Printf("[弹幕] %s：%s\n", danmaku.Sender.Uname, danmaku.Content)
		}
	})
	// 醒目留言事件
	c.OnSuperChat(func(superChat *message.SuperChat) {
		fmt.Printf("[SC|%d元] %s: %s\n", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
	})
	// 礼物事件
	c.OnGift(func(gift *message.Gift) {
		if gift.CoinType == "gold" {
			fmt.Printf("[礼物] %s 的 %s %d 个 共%.2f元\n", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
		}
	})
	// 上舰事件
	c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
		fmt.Printf("[大航海] %s 开通了 %d 等级的大航海，金额 %d 元\n", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Price/1000)
	})
	// 监听自定义事件
	c.RegisterCustomEventHandler("STOP_LIVE_ROOM_LIST", func(s string) {
		data := gjson.Get(s, "data").String()
		fmt.Printf("STOP_LIVE_ROOM_LIST: %s\n", data)
	})

	err := c.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("started")
	// 需要自行阻塞什么方法都可以
	select {}
}

func sendDanmaku() error {
	dmReq := &api.DanmakuRequest{
		Msg:      "official_13",
		RoomID:   "732",
		Bubble:   "0",
		Color:    "16777215",
		FontSize: "25",
		Mode:     "1",
		DmType:   "1",
	}
	d, err := api.SendDanmaku(dmReq, &api.BiliVerify{
		Csrf:     "",
		SessData: "",
	})
	if err != nil {
		return err
	}
	fmt.Println(d)
	return nil
}

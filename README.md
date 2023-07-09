# blivedm-go

bilibili 直播弹幕 golang 库

## 安装
```shell
go get github.com/Akegarasu/blivedm-go
```

## 快速开始

### 基础使用

该库支持以下几种基本事件，并且支持监听自定义事件。
- 弹幕
- 醒目留言
- 礼物
- 上舰
- 开播
- USER_TOAST_MSG

```go
package main

import (
	"fmt"
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	_ "github.com/Akegarasu/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func main() {
	log.SetLevel(log.DebugLevel)
	c := client.NewClient("732", "194484313") // 732为房间号，194484313为UID
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

```

### 进阶使用

#### 监听自定义事件

通过自定义监听事件，可以支持更多事件处理。  
其中，`cmd`为要监听的`cmd`名（下附常见`cmd`名）， `handler`为接收事件消息（字符串的JSON）的函数  
**注意**  
优先执行自定义 eventHandler ，会**覆盖库内自带的 handler**  
例如，如果你`RegisterCustomEventHandler("DANMU_MSG", ...`  
那么你使用`OnDanmaku`则不会再生效
```go
func (c *Client) RegisterCustomEventHandler(cmd string, handler func(s string))
```
```go
// 监听自定义事件
c.RegisterCustomEventHandler("STOP_LIVE_ROOM_LIST", func(s string) {
    data := gjson.Get(s, "data").String()
    fmt.Printf(data)
})
```

### 常见 CMD
注：来自blivedm
```python
cmd = (
        'INTERACT_WORD', 'ROOM_BANNER', 'ROOM_REAL_TIME_MESSAGE_UPDATE', 'NOTICE_MSG', 'COMBO_SEND',
        'COMBO_END', 'ENTRY_EFFECT', 'WELCOME_GUARD', 'WELCOME', 'ROOM_RANK', 'ACTIVITY_BANNER_UPDATE_V2',
        'PANEL', 'SUPER_CHAT_MESSAGE_JPN', 'USER_TOAST_MSG', 'ROOM_BLOCK_MSG', 'LIVE', 'PREPARING',
        'room_admin_entrance', 'ROOM_ADMINS', 'ROOM_CHANGE'
    )
```

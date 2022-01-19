## blivedm-go

bilibili 直播弹幕 golang 库

### 安装
```shell
go get github.com/Akegarasu/blivedm-go
```

### 快速开始

#### 基础使用

该库支持以下几种基本事件，并且支持监听自定义事件。
- 弹幕
- 醒目留言
- 礼物
- 上舰
- 开播

```go
package main

import (
	"blivedm-go/client"
	"blivedm-go/message"
	"fmt"
	"github.com/tidwall/gjson"
)

func main() {
	c := client.NewClient("8792912")
	// 弹幕事件
	c.OnDanmuku(func(danmuku *message.Danmuku) {
		fmt.Printf("[弹幕] %s：%s\n", danmuku.Sender.Uname, danmuku.Content)
	})
	// 醒目留言事件
	c.OnSuperChat(func(superChat *message.SuperChat) {
		fmt.Printf("[SC] %s: %s, %d 元\n", superChat.UserInfo.Uname, superChat.Message, superChat.Price)
	})
	// 礼物事件
	c.OnGift(func(gift *message.Gift) {
		fmt.Printf("[礼物] %s 的 %s %d 个\n", gift.Uname, gift.GiftName, gift.Num)
	})
	// 上舰事件
	c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
		fmt.Printf("%v\n", guardBuy)
	})
	err := c.ConnectAndStart()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("started")
	// 需要自行阻塞什么方法都可以
	select {}
}

```

#### 进阶使用
通过自定义监听事件，可以支持更多事件处理。
```go
// 监听自定义事件
c.RegisterCustomEventHandler("PREPARE", func(s string) {
    cmd := gjson.Get(s, "cmd").String()
    fmt.Printf(cmd)
})
```

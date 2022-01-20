package client

import (
	"bytes"
	"fmt"
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/Akegarasu/blivedm-go/packet"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"runtime/debug"
)

type eventHandlers struct {
	danmakuMessageHandlers []func(*message.Danmaku)
	superChatHandlers      []func(*message.SuperChat)
	giftHandlers           []func(*message.Gift)
	guardBuyHandlers       []func(*message.GuardBuy)
	liveHandlers           []func(*message.Live)
}

type customEventHandlers map[string]func(s string)

func (c *Client) RegisterCustomEventHandler(cmd string, handler func(s string)) {
	(*c.customEventHandlers)[cmd] = handler
}

func (c *Client) OnDanmaku(f func(*message.Danmaku)) {
	c.eventHandlers.danmakuMessageHandlers = append(c.eventHandlers.danmakuMessageHandlers, f)
}

func (c *Client) OnSuperChat(f func(*message.SuperChat)) {
	c.eventHandlers.superChatHandlers = append(c.eventHandlers.superChatHandlers, f)
}

func (c *Client) OnGift(f func(gift *message.Gift)) {
	c.eventHandlers.giftHandlers = append(c.eventHandlers.giftHandlers, f)
}

func (c *Client) OnGuardBuy(f func(*message.GuardBuy)) {
	c.eventHandlers.guardBuyHandlers = append(c.eventHandlers.guardBuyHandlers, f)
}

func (c *Client) OnLive(f func(*message.Live)) {
	c.eventHandlers.liveHandlers = append(c.eventHandlers.liveHandlers, f)
}

func (c *Client) Handle(p packet.Packet) {
	switch p.Operation {
	case packet.Notification:
		sb := bytes.NewBuffer(p.Body).String()
		cmd := gjson.Get(sb, "cmd").String()
		// 优先执行自定义 eventHandler ，会覆盖库内自带的 handler
		f, ok := (*c.customEventHandlers)[cmd]
		if ok {
			go cover(func() { f(sb) })
			return
		}
		switch cmd {
		case "DANMU_MSG":
			d := new(message.Danmaku)
			d.Parse(p.Body)
			for _, fn := range c.eventHandlers.danmakuMessageHandlers {
				go cover(func() { fn(d) })
			}
		case "SUPER_CHAT_MESSAGE":
			s := new(message.SuperChat)
			s.Parse(p.Body)
			for _, fn := range c.eventHandlers.superChatHandlers {
				go cover(func() { fn(s) })
			}
		case "SEND_GIFT":
			g := new(message.Gift)
			g.Parse(p.Body)
			for _, fn := range c.eventHandlers.giftHandlers {
				go cover(func() { fn(g) })
			}
		case "GUARD_BUY":
			g := new(message.GuardBuy)
			g.Parse(p.Body)
			for _, fn := range c.eventHandlers.guardBuyHandlers {
				go cover(func() { fn(g) })
			}
		case "LIVE":
			l := new(message.Live)
			l.Parse(p.Body)
			for _, fn := range c.eventHandlers.liveHandlers {
				go cover(func() { fn(l) })
			}
		//TODO: cmd补全
		case "INTERACT_WORD":
		case "ROOM_BANNER":
		case "ROOM_REAL_TIME_MESSAGE_UPDATE":
		case "NOTICE_MSG":
		case "COMBO_SEND":
		case "COMBO_END":
		case "ENTRY_EFFECT":
		case "WELCOME_GUARD":
		case "WELCOME":
		case "ROOM_RANK":
		case "ACTIVITY_BANNER_UPDATE_V2":
		case "PANEL":
		case "SUPER_CHAT_MESSAGE_JPN":
		case "USER_TOAST_MSG":
		case "ROOM_BLOCK_MSG":
		case "PREPARING":
		case "room_admin_entrance":
		case "ROOM_ADMINS":
		case "ROOM_CHANGE":
		case "LIVE_INTERACTIVE_GAME":
		case "WIDGET_BANNER":
		case "ONLINE_RANK_COUNT":
		case "ONLINE_RANK_V2":
		case "STOP_LIVE_ROOM_LIST":
		case "ONLINE_RANK_TOP3":
		case "HOT_RANK_CHANGED":
		case "HOT_RANK_CHANGED_V2":
		default:
			//log.Infof("cmd %s, %s", p.Body, cmd)
			log.WithField("data", string(p.Body)).Warn("unknown cmd")
		}
	case packet.HeartBeatResponse:
	case packet.RoomEnterResponse:
	default:
		log.WithField("protover", p.ProtocolVersion).
			WithField("data", string(p.Body)).
			Warn("unknown protover")
	}
}

func cover(f func()) {
	defer func() {
		if pan := recover(); pan != nil {
			fmt.Printf("event error: %v\n%s", pan, debug.Stack())
		}
	}()
	f()
}

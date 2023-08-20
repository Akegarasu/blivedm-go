package client

import (
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/Akegarasu/blivedm-go/packet"
	"github.com/Akegarasu/blivedm-go/utils"
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"strings"
)

var (
	knownCMD    = []string{"INTERACT_WORD", "HOT_RANK_SETTLEMENT", "DANMU_GIFT_LOTTERY_START", "WELCOME_GUARD", "PK_PROCESS", "PK_BATTLE_PRO_TYPE", "MATCH_TEAM_GIFT_RANK", "PK_BATTLE_CRIT", "LUCK_GIFT_AWARD_USER", "SCORE_CARD", "ONLINE_RANK_V2", "PK_BATTLE_SPECIAL_GIFT", "SEND_TOP", "SUPER_CHAT_MESSAGE_JPN", "ANIMATION", "GUARD_LOTTERY_START", "WEEK_STAR_CLOCK", "WELCOME", "WIN_ACTIVITY", "ROOM_KICKOUT", "CHANGE_ROOM_INFO", "ROOM_SKIN_MSG", "ROOM_BLOCK_MSG", "SUPER_CHAT_ENTRANCE", "PK_BATTLE_RANK_CHANGE", "ROOM_LOCK", "TV_END", "PK_PRE", "ROOM_SILENT_OFF", "SEND_GIFT", "DANMU_MSG", "ANCHOR_LOT_START", "ROOM_BOX_USER", "ONLINE_RANK_TOP3", "WIDGET_BANNER", "PK_BATTLE_START", "ACTIVITY_MATCH_GIFT", "PK_AGAIN", "PK_MATCH", "RAFFLE_START", "LIVE", "WISH_BOTTLE", "GUARD_ACHIEVEMENT_ROOM", "ONLINE_RANK_COUNT", "COMMON_NOTICE_DANMAKU", "LOL_ACTIVITY", "HOT_RANK_CHANGED", "ROOM_BLOCK_INTO", "ROOM_LIMIT", "PANEL", "RAFFLE_END", "ENTRY_EFFECT", "STOP_LIVE_ROOM_LIST", "TV_START", "WATCH_LPL_EXPIRED", "PK_BATTLE_PRE", "USER_TOAST_MSG", "BOX_ACTIVITY_START", "PK_MIC_END", "LIVE_INTERACTIVE_GAME", "ROOM_BANNER", "PK_BATTLE_GIFT", "MESSAGEBOX_USER_GAIN_MEDAL", "LITTLE_TIPS", "HOUR_RANK_AWARDS", "NOTICE_MSG", "ROOM_REAL_TIME_MESSAGE_UPDATE", "ANCHOR_LOT_END", "PREPARING", "GUARD_BUY", "ROOM_CHANGE", "room_admin_entrance", "CHASE_FRAME_SWITCH", "DANMU_GIFT_LOTTERY_AWARD", "PK_BATTLE_VOTES_ADD", "PK_BATTLE_END", "CUT_OFF", "PK_BATTLE_PROCESS", "PK_BATTLE_SETTLE_USER", "ANCHOR_LOT_AWARD", "WIN_ACTIVITY_USER", "VOICE_JOIN_STATUS", "DANMU_GIFT_LOTTERY_END", "ROOM_RANK", "SUPER_CHAT_MESSAGE", "ACTIVITY_BANNER_UPDATE_V2", "SPECIAL_GIFT", "ROOM_SILENT_ON", "WARNING", "ROOM_ADMINS", "COMBO_SEND", "HOT_RANK_SETTLEMENT_V2", "ANCHOR_LOT_CHECKSTATUS", "HOT_RANK_CHANGED_V2", "SUPER_CHAT_MESSAGE_DELETE", "PK_END", "PK_SETTLE", "ROOM_REFRESH", "PK_START", "COMBO_END", "PK_LOTTERY_START", "GUARD_WINDOWS_OPEN", "REENTER_LIVE_ROOM", "MESSAGEBOX_USER_MEDAL_CHANGE", "MESSAGEBOX_USER_MEDAL_COMPENSATION", "LITTLE_MESSAGE_BOX", "PK_BATTLE_PRE_NEW", "PK_BATTLE_START_NEW", "PK_BATTLE_PROCESS_NEW", "PK_BATTLE_FINAL_PROCESS", "PK_BATTLE_SETTLE_V2", "PK_BATTLE_SETTLE_NEW", "PK_BATTLE_PUNISH_END", "PK_BATTLE_VIDEO_PUNISH_BEGIN", "PK_BATTLE_VIDEO_PUNISH_END", "ENTRY_EFFECT_MUST_RECEIVE", "SUPER_CHAT_AUDIT", "VIDEO_CONNECTION_JOIN_START", "VIDEO_CONNECTION_JOIN_END", "VIDEO_CONNECTION_MSG", "VTR_GIFT_LOTTERY", "RED_POCKET_START", "FULL_SCREEN_SPECIAL_EFFECT", "POPULARITY_RED_POCKET_START", "POPULARITY_RED_POCKET_WINNER_LIST", "USER_PANEL_RED_ALARM", "SHOPPING_CART_SHOW", "THERMAL_STORM_DANMU_BEGIN", "THERMAL_STORM_DANMU_UPDATE", "THERMAL_STORM_DANMU_CANCEL", "THERMAL_STORM_DANMU_OVER", "MILESTONE_UPDATE_EVENT", "WEB_REPORT_CONTROL", "DANMU_TAG_CHANGE", "RANK_REM", "LIVE_PLAYER_LOG_RECYCLE", "LIVE_INTERNAL_ROOM_LOGIN", "LIVE_OPEN_PLATFORM_GAME", "WATCHED_CHANGE", "DANMU_AGGREGATION", "POPULARITY_RED_POCKET_NEW", "LIKE_INFO_V3_CLICK", "POPULAR_RANK_CHANGED"}
	knownCMDMap map[string]int
)

type eventHandlers struct {
	danmakuMessageHandlers []func(*message.Danmaku)
	superChatHandlers      []func(*message.SuperChat)
	giftHandlers           []func(*message.Gift)
	guardBuyHandlers       []func(*message.GuardBuy)
	liveHandlers           []func(*message.Live)
	userToastHandlers      []func(*message.UserToast)
}

type customEventHandlers map[string]func(s string)

func init() {
	knownCMDMap = make(map[string]int)
	for _, c := range knownCMD {
		knownCMDMap[c] = 0
	}
}

// RegisterCustomEventHandler 注册 自定义事件 的处理器
//
// 需要提供事件名，可参考 knownCMD
func (c *Client) RegisterCustomEventHandler(cmd string, handler func(s string)) {
	(*c.customEventHandlers)[cmd] = handler
}

// OnDanmaku 添加 弹幕事件 的处理器
func (c *Client) OnDanmaku(f func(*message.Danmaku)) {
	c.eventHandlers.danmakuMessageHandlers = append(c.eventHandlers.danmakuMessageHandlers, f)
}

// OnSuperChat 添加 醒目留言事件 的处理器
func (c *Client) OnSuperChat(f func(*message.SuperChat)) {
	c.eventHandlers.superChatHandlers = append(c.eventHandlers.superChatHandlers, f)
}

// OnGift 添加 礼物事件 的处理器
func (c *Client) OnGift(f func(gift *message.Gift)) {
	c.eventHandlers.giftHandlers = append(c.eventHandlers.giftHandlers, f)
}

// OnGuardBuy 添加 开通大航海事件 的处理器
func (c *Client) OnGuardBuy(f func(*message.GuardBuy)) {
	c.eventHandlers.guardBuyHandlers = append(c.eventHandlers.guardBuyHandlers, f)
}

// OnLive 添加 开播事件 的处理器
func (c *Client) OnLive(f func(*message.Live)) {
	c.eventHandlers.liveHandlers = append(c.eventHandlers.liveHandlers, f)
}

// OnUserToast 添加 UserToast 的处理器
func (c *Client) OnUserToast(f func(*message.UserToast)) {
	c.eventHandlers.userToastHandlers = append(c.eventHandlers.userToastHandlers, f)
}

// Handle 处理一个包
func (c *Client) Handle(p packet.Packet) {
	switch p.Operation {
	case packet.Notification:
		cmd := parseCmd(p.Body)
		sb := utils.BytesToString(p.Body)
		// 新的弹幕 cmd 可能带参数
		if ind := strings.Index(cmd, ":"); ind >= 0 {
			cmd = cmd[:ind]
		}
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
		case "USER_TOAST_MSG":
			u := new(message.UserToast)
			u.Parse(p.Body)
			for _, fn := range c.eventHandlers.userToastHandlers {
				go cover(func() { fn(u) })
			}
		default:
			if _, ok := knownCMDMap[cmd]; ok {
				return
			}
			log.Debugf("unknown cmd(%s), body: %s", cmd, p.Body)
		}
	case packet.HeartBeatResponse:
	case packet.RoomEnterResponse:
	default:
		log.WithField("protover", p.ProtocolVersion).
			WithField("data", string(p.Body)).
			Warn("unknown protover")
	}
}

// parseCmd 获取 JSON 报文的 CMD
func parseCmd(d []byte) string {
	// {"cmd":"DANMU_MSG", ...
	l := len(d)
	pos := 8
	s := byte('"')
	for pos <= l {
		if d[pos] == s {
			break
		}
		pos++
	}
	return utils.BytesToString(d[8:pos])
}

func cover(f func()) {
	defer func() {
		if pan := recover(); pan != nil {
			log.Errorf("event error: %v\n%s", pan, debug.Stack())
		}
	}()
	f()
}

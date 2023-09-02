package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/Akegarasu/blivedm-go/api"
	"github.com/Akegarasu/blivedm-go/packet"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	conn                *websocket.Conn
	Uid                 string
	Buvid               string
	RoomID              string
	Cookie              string
	token               string
	host                string
	hostList            []string
	retryCount          int
	eventHandlers       *eventHandlers
	customEventHandlers *customEventHandlers
	cancel              context.CancelFunc
	done                <-chan struct{}
}

// NewClient 创建一个新的弹幕 client
func NewClient(roomID string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		RoomID:              roomID,
		retryCount:          0,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
		done:                ctx.Done(),
		cancel:              cancel,
	}
}

func (c *Client) SetCookie(cookie string) {
	c.Cookie = cookie
}

// init 初始化 获取真实 RoomID 和 弹幕服务器 host
func (c *Client) init() error {
	if c.Cookie != "" {
		uid, err := api.GetUid(c.Cookie)
		if err != nil {
			if c.Uid == "" {
				c.Uid = "0"
			}
		}
		c.Uid = uid
		re := regexp.MustCompile("_uuid=(.+?);")
		result := re.FindAllStringSubmatch(c.Cookie, -1)
		if len(result) > 0 {
			c.Buvid = result[0][1]
		}
	}
	roomInfo, err := api.GetRoomInfo(c.RoomID)
	// 失败降级
	if err != nil || roomInfo.Code != 0 {
		log.Errorf("room=%s init GetRoomInfo fialed, %s", c.RoomID, err)
	}
	c.RoomID = strconv.Itoa(roomInfo.Data.RoomId)
	if c.host == "" {
		info, err := api.GetDanmuInfo(c.RoomID, c.Cookie)
		if err != nil {
			c.hostList = []string{"broadcastlv.chat.bilibili.com"}
		} else {
			for _, h := range info.Data.HostList {
				c.hostList = append(c.hostList, h.Host)
			}
		}
		c.token = info.Data.Token
	}
	return nil
}

func (c *Client) connect() error {
	reqHeader := &http.Header{}
	reqHeader.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
retry:
	c.host = c.hostList[c.retryCount%len(c.hostList)]
	c.retryCount++
	conn, res, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s/sub", c.host), *reqHeader)
	if err != nil {
		log.Errorf("connect dial failed, retry %d times", c.retryCount)
		time.Sleep(2 * time.Second)
		goto retry
	}
	c.conn = conn
	_ = res.Body.Close()
	if err = c.sendEnterPacket(); err != nil {
		log.Errorf("failed to send enter packet, retry %d times", c.retryCount)
		time.Sleep(2 * time.Second)
		goto retry
	}
	return nil
}

func (c *Client) wsLoop() {
	for {
		select {
		case <-c.done:
			log.Debug("current client closed")
			return
		default:
			msgType, data, err := c.conn.ReadMessage()
			if err != nil {
				log.Error("ws message read failed, reconnecting")
				time.Sleep(time.Duration(3) * time.Millisecond)
				_ = c.connect()
				continue
			}
			if msgType != websocket.BinaryMessage {
				log.Error("packet not binary")
				continue
			}
			for _, pkt := range packet.DecodePacket(data).Parse() {
				go c.Handle(pkt)
			}
		}
	}
}

func (c *Client) heartBeatLoop() {
	pkt := packet.NewHeartBeatPacket()
	for {
		select {
		case <-c.done:
			return
		case <-time.After(30 * time.Second):
			if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
				log.Error(err)
			}
			log.Debug("send: HeartBeat")
		}
	}
}

// Start 启动弹幕 Client 初始化并连接 ws、发送心跳包
func (c *Client) Start() error {
	if err := c.init(); err != nil {
		return err
	}
	if err := c.connect(); err != nil {
		return err
	}
	go c.wsLoop()
	go c.heartBeatLoop()
	return nil
}

// Stop 停止弹幕 Client
func (c *Client) Stop() {
	c.cancel()
}

func (c *Client) SetHost(host string) {
	c.host = host
}

// UseDefaultHost 使用默认 host broadcastlv.chat.bilibili.com
func (c *Client) UseDefaultHost() {
	c.hostList = []string{"broadcastlv.chat.bilibili.com"}
}

func (c *Client) sendEnterPacket() error {
	rid, err := strconv.Atoi(c.RoomID)
	if err != nil {
		return errors.New("error RoomID")
	}
	uid, err := strconv.Atoi(c.Uid)
	if err != nil {
		return errors.New("error UID")
	}
	pkt := packet.NewEnterPacket(uid, c.Buvid, rid, c.token)
	if err = c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		return err
	}
	log.Debugf("send: EnterPacket")
	return nil
}

package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Akegarasu/blivedm-go/api"
	"github.com/Akegarasu/blivedm-go/packet"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	conn                *websocket.Conn
	uid                 int
	roomID              string
	tempID              string
	token               string
	host                string
	hostList            []string
	eventHandlers       *eventHandlers
	customEventHandlers *customEventHandlers
	cancel              context.CancelFunc
	done                <-chan struct{}
}

// NewClient 创建一个新的弹幕 client
func NewClient(roomID string, uid int) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		tempID:              roomID,
		uid:                 uid,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
		done:                ctx.Done(),
		cancel:              cancel,
	}
}

// init 初始化 获取真实 roomID 和 弹幕服务器 host
func (c *Client) init() error {
	rid, _ := strconv.Atoi(c.tempID)
	// 处理 shortID
	if rid <= 1000 && c.roomID == "" {
		realID, err := api.GetRoomRealID(c.tempID)
		if err != nil {
			return err
		}
		c.roomID = realID
	} else {
		c.roomID = c.tempID
	}
	if c.host == "" {
		info, err := api.GetDanmuInfo(c.roomID)
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
	retryCount := 0
retry:
	// 随着重连会自动切换弹幕服务器
	c.host = c.hostList[retryCount%len(c.hostList)]
	retryCount++
	conn, res, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s/sub", c.host), nil)
	if err != nil {
		log.Errorf("connect dial failed, retry %d times", retryCount)
		time.Sleep(2 * time.Second)
		goto retry
	}
	c.conn = conn
	res.Body.Close()
	if err = c.sendEnterPacket(); err != nil {
		log.Errorf("failed to send enter packet, retry %d times", retryCount)
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
				log.Info("reconnect")
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
	rid, err := strconv.Atoi(c.roomID)
	if err != nil {
		return errors.New("error roomID")
	}
	pkt := packet.NewEnterPacket(c.uid, rid, c.token)
	if err = c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		return err
	}
	log.Debugf("send: EnterPacket")
	return nil
}

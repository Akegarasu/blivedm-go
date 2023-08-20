package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Akegarasu/blivedm-go/api"
	"github.com/Akegarasu/blivedm-go/packet"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	conn                *websocket.Conn
	uid                 string
	roomID              string
	tempID              string
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
func NewClient(roomID string, uid string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		tempID:              roomID,
		uid:                 uid,
		retryCount:          0,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
		done:                ctx.Done(),
		cancel:              cancel,
	}
}

// init 初始化 获取真实 roomID 和 弹幕服务器 host
func (c *Client) init() error {
	roomInfo, err := api.GetRoomInfo(c.tempID)
	// 失败降级
	if err != nil || roomInfo.Code != 0 {
		log.Errorf("room=%s init GetRoomInfo fialed, %s", c.tempID, err)
		c.roomID = c.tempID
		if c.uid == "" {
			c.uid = "208259" // 叔叔的 UID
		}
	}
	c.roomID = strconv.Itoa(roomInfo.Data.RoomId)
	if c.uid == "" {
		c.uid = strconv.Itoa(roomInfo.Data.Uid)
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
	res.Body.Close()
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
	rid, err := strconv.Atoi(c.roomID)
	if err != nil {
		return errors.New("error roomID")
	}
	uid, err := strconv.Atoi(c.uid)
	if err != nil {
		return errors.New("error UID")
	}
	pkt := packet.NewEnterPacket(uid, rid, c.token)
	if err = c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		return err
	}
	log.Debugf("send: EnterPacket")
	return nil
}

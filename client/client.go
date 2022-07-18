package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/Akegarasu/blivedm-go/api"
	"github.com/Akegarasu/blivedm-go/packet"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Client struct {
	conn                *websocket.Conn
	roomID              string
	realRoomID          string
	token               string
	host                string
	hostList            []string
	eventHandlers       *eventHandlers
	customEventHandlers *customEventHandlers
	cancel              context.CancelFunc
	done                <-chan struct{}
}

func NewClient(roomID string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		roomID:              roomID,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
		done:                ctx.Done(),
		cancel:              cancel,
	}
}

func (c *Client) Connect() error {
	retryCount := 0
	rid, _ := strconv.Atoi(c.roomID)
	if rid <= 1000 && c.realRoomID == "" {
		realID, err := api.GetRoomRealID(c.roomID)
		if err != nil {
			return err
		}
		c.roomID = realID
		c.realRoomID = realID
	} else {
		c.realRoomID = c.roomID
	}
	if c.host == "" {
		info, err := api.GetDanmuInfo(c.realRoomID)
		if err != nil {
			c.hostList = []string{"broadcastlv.chat.bilibili.com"}
		} else {
			for _, h := range info.Data.HostList {
				c.hostList = append(c.hostList, h.Host)
			}
		}
		c.token = info.Data.Token
	}
retry:
	c.host = c.hostList[retryCount%len(c.hostList)]
	retryCount++
	conn, res, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s/sub", c.host), nil)
	if err != nil {
		log.Error("connect dial failed, retry...")
		time.Sleep(2 * time.Second)
		goto retry
	}
	c.conn = conn
	res.Body.Close()
	if err = c.sendEnterPacket(); err != nil {
		log.Error("connect enter packet send failed, retry...")
		goto retry
	}
	return nil
}

func (c *Client) listen() {
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
				_ = c.Connect()
				continue
			}
			if msgType != websocket.BinaryMessage {
				log.Error("packet not binary", data)
				continue
			}
			for _, pkt := range packet.DecodePacket(data).Parse() {
				go c.Handle(pkt)
			}
		}
	}
}

func (c *Client) Start() {
	go c.listen()
	go c.heartBeat()
}

func (c *Client) Stop() {
	c.cancel()
}

func (c *Client) ConnectAndStart() error {
	if err := c.Connect(); err != nil {
		return err
	}
	c.Start()
	return nil
}

func (c *Client) SetHost(host string) {
	c.host = host
}

// UseDefaultHost 使用默认 host broadcastlv.chat.bilibili.com
func (c *Client) UseDefaultHost() {
	c.SetHost("broadcastlv.chat.bilibili.com")
}

func (c *Client) heartBeat() {
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

func (c *Client) sendEnterPacket() error {
	rid, err := strconv.Atoi(c.realRoomID)
	if err != nil {
		return errors.New("error roomID")
	}
	pkt := packet.NewEnterPacket(0, rid, c.token)
	if err = c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		return err
	}
	log.Debugf("send: EnterPacket")
	return nil
}

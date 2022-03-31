package client

import (
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
	token               string
	host                string
	eventHandlers       *eventHandlers
	customEventHandlers *customEventHandlers
	done                chan struct{}
}

func NewClient(roomID string) *Client {
	return &Client{
		roomID:              roomID,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
		done:                make(chan struct{}),
	}
}

func (c *Client) Connect() error {
	rid, _ := strconv.Atoi(c.roomID)
	if rid <= 1000 {
		realID, err := api.GetRoomRealID(c.roomID)
		if err != nil {
			return err
		}
		c.roomID = realID
	}
	if c.host == "" {
		info, err := api.GetDanmuInfo(c.roomID)
		if err != nil {
			return err
		}
		c.host = fmt.Sprintf("wss://%s/sub", info.Data.HostList[0].Host)
		c.token = info.Data.Token
	}
	conn, _, err := websocket.DefaultDialer.Dial(c.host, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Start() error {
	if err := c.sendEnterPacket(); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-c.done:
				return
			default:
				msgType, data, err := c.conn.ReadMessage()
				if err != nil {
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
	}()
	go c.startHeartBeat()
	return nil
}

func (c *Client) Stop() {
	close(c.done)
}

func (c *Client) ConnectAndStart() error {
	if err := c.Connect(); err != nil {
		return err
	}
	if err := c.Start(); err != nil {
		return err
	}
	return nil
}

func (c *Client) SetHost(host string) {
	c.host = host
}

func (c *Client) UseDefaultHost() {
	c.SetHost("wss://broadcastlv.chat.bilibili.com/sub")
}

func (c *Client) startHeartBeat() {
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
	rid, err := strconv.Atoi(c.roomID)
	if err != nil {
		return errors.New("error roomID")
	}
	pkt := packet.NewEnterPacket(0, rid)
	if err = c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		return err
	}
	log.Debugf("send: EnterPacket: %v", pkt)
	return nil
}

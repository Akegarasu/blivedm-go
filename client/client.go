package client

import (
	"encoding/json"
	"fmt"
	"github.com/Akegarasu/blivedm-go/packet"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
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
}

func NewClient(roomID string) *Client {
	return &Client{
		roomID:              roomID,
		eventHandlers:       &eventHandlers{},
		customEventHandlers: &customEventHandlers{},
	}
}

func (c *Client) Connect() error {
	if c.host == "" {
		info, err := getDanmuInfo(c.roomID)
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

func (c *Client) Start() {
	c.sendEnterPacket()
	go func() {
		for {
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
	}()
	go c.startHeartBeat()
}

func (c *Client) ConnectAndStart() error {
	err := c.Connect()
	if err != nil {
		return err
	}
	c.Start()
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
		if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
			log.Fatal(err)
		}
		log.Debug("send: HeartBeat")
		time.Sleep(30 * time.Second)
	}
}

func (c *Client) sendEnterPacket() {
	rid, err := strconv.Atoi(c.roomID)
	if err != nil {
		log.Fatal("error roomID")
	}
	pkt := packet.NewEnterPacket(0, rid)
	if err := c.conn.WriteMessage(websocket.BinaryMessage, pkt); err != nil {
		log.Fatal(err)
	}
	log.Debugf("send: EnterPacket: %v", pkt)
}

func getDanmuInfo(roomID string) (*DanmuInfo, error) {
	url := fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%s&type=0", roomID)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	result := &DanmuInfo{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}

package sender

import (
	"XAXAtonMTC/pkg/packetsender"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Sender struct {
	wsConn *websocket.Conn
}

func NewSender(w http.ResponseWriter, r *http.Request) (*Sender, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &Sender{
		wsConn: conn,
	}, nil
}

func (sender *Sender) SendPacket(packet *packetsender.Packet) error {
	if sender.wsConn == nil {
		return NoWebsocketConnectionErr
	}

	w, err := sender.wsConn.NextWriter(websocket.BinaryMessage)
	defer w.Close()
	if err != nil {
		return err
	}

	buf, err := json.Marshal(*packet)
	if err != nil {
		return err
	}

	var pack packetsender.Packet
	json.Unmarshal(buf, &pack)

	f, _ := os.OpenFile("e.mp3", os.O_APPEND, 0640)
	f.Write(pack.Data)
	f.Close()

	_, err = w.Write(buf)

	return err
}

type SendersRepo interface {
	AddSender(uint32, http.ResponseWriter, *http.Request) error
	SendPacket([]uint32, func(uint32) bool, *packetsender.Packet)
	DeleteSender(uint32)
	CloseConnections(...uint32)
}

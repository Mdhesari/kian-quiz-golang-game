package websocketadapter

import (
	"bytes"
	"encoding/json"
	"mdhesari/kian-quiz-golang-game/logger"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	cfg  *ClientConfig
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type ClientConfig struct {
	PingPeriod     time.Duration `koanf:"ping_period"`
	ReadTimeout    time.Duration `koanf:"read_timeout"`
	WriteTimeout   time.Duration `koanf:"write_timeout"`
	MaxMessageSize int64         `koanf:"max_message_size"`
}

func (a Adapter) NewClient(conn *websocket.Conn) *Client {
	return &Client{
		cfg:  &a.cfg.ClientCfg,
		hub:  a.hub,
		conn: conn,
		send: make(chan []byte),
	}
}

func (c Client) ReadPump() {
	defer func() {
		c.hub.unregister <- &c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(c.cfg.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.cfg.ReadTimeout))
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(c.cfg.ReadTimeout)); err != nil {
			return err
		}
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.L().Error("Error reading from websocket.", zap.Error(err))
			}

			break
		}

		var msg struct {
			Type    string      `json:"type"`
			Payload interface{} `json:"payload"`
		}

		if err = json.Unmarshal(message, msg); err != nil {
			logger.L().Error("Could not unmarshal message.", zap.Error(err))

			break
		}

		logger.L().Info("new msg from cli", zap.Any("msg", message))

		// c.hub.broadcast <- message
	}
}

func (c Client) WritePump() {
	ticker := time.NewTicker(c.cfg.PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.L().Error("Error writing to websocket.", zap.Error(err))

				return
			}
			if _, err := w.Write(message); err != nil {
				logger.L().Error("Error writing to websocket.", zap.Error(err))

				w.Close()

				return
			}

			if err := w.Close(); err != nil {
				logger.L().Error("Error closing writer.", zap.Error(err))
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
				logger.L().Error("Error sending ping.", zap.Error(err))

				return
			}

			logger.L().Info("ping")
		}
	}
}
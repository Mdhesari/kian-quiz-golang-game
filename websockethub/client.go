package websockethub

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/protobufdecoder"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Client struct {
	conn   *websocket.Conn
	hub    *Hub
	userID primitive.ObjectID
	send   chan []byte
	stop   chan struct{}
}

func NewClient(conn *websocket.Conn, hub *Hub, userID primitive.ObjectID) Client {
	return Client{
		conn:   conn,
		hub:    hub,
		userID: userID,
		send:   make(chan []byte),
		stop:   make(chan struct{}),
	}
}

func (c *Client) Start() {
	logger.L().Info("Starting client writes and reads.")

	go c.writePump()
	go c.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.stop <- struct{}{}

		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			logger.L().Error("Pong handler: read deadline err.", zap.Error(err))
		}

		logger.L().Info("pong")

		return nil
	})

	for {
		logger.L().Info("going through messages.")
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.L().Error("Could not read messages from websocket client.", zap.Error(err))
			}

			return
		}

		websocketMsg := protobufdecoder.DecodeWebSocketMsg(string(msg))

		logger.L().Info("message successfully recieved.", zap.Any("message", websocketMsg.Payload))

		c.handleWebsocketMsg(websocketMsg)
	}
}

func (c *Client) handleWebsocketMsg(msg entity.WebsocketMsg) {
	switch msg.Type {
	case string(entity.GameStartedEvent):
		playersMatched := protobufdecoder.DecodeGameStartedEvent(msg.Payload)

		logger.L().Info("Game started.", zap.Any("playersMatched", playersMatched))
	case string(entity.PlayerAnsweredEvent):
		c.hub.Publish(context.Background(), string(entity.PlayerAnsweredEvent), msg.Payload)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.stop <- struct{}{}

		c.conn.Close()
	}()

	ticker := time.NewTicker(pingPeriod)

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				logger.L().Error("Pong handler: write deadline err")

				c.conn.WriteMessage(websocket.CloseMessage, []byte{})

				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			logger.L().Info("Writing message.", zap.Any("message", msg))
			if _, err := w.Write(msg); err != nil {
				logger.L().Error("Could not write msg.", zap.Any("msg", msg))
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-c.stop:
			ticker.Stop()

			return
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.L().Error("Could not ping the client.", zap.Error(err))

				return
			}
		}
	}
}

func (c *Client) cleanUp() {
	close(c.send)
	c.conn.Close()
}

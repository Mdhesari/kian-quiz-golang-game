package websockethub

import (
	"encoding/json"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

type Client struct {
	conn   *websocket.Conn
	userID primitive.ObjectID
	send   chan []byte
	stop   chan struct{}
}

type Message struct {
	Type   string
	UserID string
	Body   []byte
}

func NewHub() Hub {
	return Hub{
		clients:    map[string]*Client{},
		register:   make(chan *Client, 1000),
		unregister: make(chan *Client, 1000),
		broadcast:  make(chan *Message, 1000),
	}
}

func (h *Hub) RegisterClient(cli *Client) {
	select {
	case h.register <- cli:
	default:
		logger.L().Warn("Could not register client.", zap.Any("client", cli))
	}
}

func (h *Hub) BroadcastMessage(msg *Message) {
	if msg == nil {
		logger.L().Error("Attempted to broadcast a nil message")
		return
	}

	select {
	case h.broadcast <- msg:
	default:
		logger.L().Warn("Broadcast is full.")
	}

	logger.L().Info("done broadcast.")
}

func (h *Hub) Start() {
	for {
		select {
		case cli := <-h.register:
			logger.L().Info("Registering cli", zap.Any("cli", cli))

			h.clients[cli.userID.Hex()] = cli

			logger.L().Info("Registered cli")
		case cli := <-h.unregister:
			logger.L().Info("Unregistering cli", zap.Any("cli", cli))

			delete(h.clients, cli.userID.Hex())
			cli.cleanUp()

			logger.L().Info("done unregistering cli")
		case msg := <-h.broadcast:
			finalMsg := protobufencoder.EncodeWebSocketMsg(entity.WebsocketMsg{
				Type:    msg.Type,
				Payload: msg.Body,
			})

			if msg.UserID != "" {
				cli, ok := h.clients[msg.UserID]
				if !ok {
					logger.L().Error("Client is not available.", zap.String("userId", msg.UserID), zap.Any("clients", h.clients))

					continue
				}

				select {
				case cli.send <- []byte(finalMsg):
				default:
					close(cli.send)
					delete(h.clients, cli.userID.Hex())
				}

				continue
			}

			for userID, cli := range h.clients {
				select {
				case cli.send <- []byte(finalMsg):
				default:
					close(cli.send)
					delete(h.clients, userID)
				}
			}
		}
	}
}

func NewClient(conn *websocket.Conn, userID primitive.ObjectID) Client {
	return Client{
		conn:   conn,
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

		logger.L().Info("decoded", zap.Any("message", string(msg)))
		var websocketMsg entity.WebsocketMsg
		if err := json.Unmarshal(msg, &websocketMsg); err != nil {
			logger.L().Error("Could not decode websocket msg.", zap.Error(err))

			return
		}

		logger.L().Info("message successfully recieved.", zap.Any("message", websocketMsg))

		// TODO - Handle cli messages (game answer, game ending, etc)
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
			w.Write(msg)

			if err := w.Close(); err != nil {
				return
			}
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

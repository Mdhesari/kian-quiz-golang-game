package websockethub

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
	"time"

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

type Publsiher interface {
	Publish(ctx context.Context, topic string, payload string)
}

type Hub struct {
	pub        Publsiher
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

type Message struct {
	Type    string
	UserIDs []string
	Body    string
}

func NewHub(pub Publsiher) Hub {
	return Hub{
		pub:        pub,
		clients:    map[string]*Client{},
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (h *Hub) Publish(ctx context.Context, topic string, payload string) {
	h.pub.Publish(ctx, topic, payload)
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

	h.broadcast <- msg

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

			var clients map[string]*Client = make(map[string]*Client)

			if len(msg.UserIDs) > 0 {
				for _, userID := range msg.UserIDs {
					var ok bool
					userId, ok := h.clients[userID]
					if !ok {
						// TODO - Send msg to queue
						logger.L().Warn("Could not find user id in clients.", zap.Any("userID", userId))
					} else {
						clients[userID] = userId
					}
				}

				logger.L().Info("Broadcasting to specific clients.", zap.Any("userIDs", msg.UserIDs))
			} else {
				clients = h.clients

				logger.L().Info("Broadcasting to all clients.")
			}

			logger.L().Info("send to clients.", zap.Any("clients", clients))

			for userID, cli := range clients {
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

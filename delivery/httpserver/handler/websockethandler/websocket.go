package websockethandler

import (
	"context"
	"encoding/json"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/claim"
	"mdhesari/kian-quiz-golang-game/pkg/timestamp"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (h Handler) Websocket(c echo.Context) error {
	claims := claim.GetClaimsFromEchoContext(c)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO - check origin
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		logger.L().Error("Could not upgrade http to websocket.", zap.Error(err))

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var ch <-chan *redis.Message

	go func() {
		ch = h.pubsubManager.SubscribeAndGetChannel(string(entity.GameStartedEvent))
	}()

	go func(conn *websocket.Conn, ch <-chan *redis.Message) {
		defer conn.Close()

		ticker := time.NewTicker(60 * time.Second)
		defer func() {
			ticker.Stop()
			conn.Close()
		}()

		for {
			select {
			case msg := <-ch:
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
					logger.L().Error("Error writing message to websocket.", zap.Error(err), zap.String("msg", msg.Payload))

					return
				}
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
					logger.L().Error("Error sending ping.", zap.Error(err))

					return
				}

				logger.L().Info("ping")
			}
		}
	}(conn, ch)

	go func(conn *websocket.Conn) {
		defer conn.Close()

		conn.SetPongHandler(func(msg string) error {
			if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
				return err
			}

			go func() {
				userId := claims.UserID
				_, err := h.presenceSrv.Upsert(context.Background(), param.PresenceUpsertRequest{
					UserId:    userId,
					Timestamp: timestamp.Now(),
				})
				if err != nil {
					logger.L().Error("Error upserting presence after pong.", zap.Error(err))
				}
			}()

			return nil
		})

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.L().Error("Error reading from websocket.", zap.Error(err))
				}

				break
			}
			logger.L().Info("new msg from cli", zap.Any("msg", msg))

			var decodedMsg entity.WebsocketMsg
			if err := json.Unmarshal(msg, &decodedMsg); err != nil {
				logger.L().Error("Could not unmarshal cli message.", zap.Error(err))
			}

			switch decodedMsg.Type {
			case "answer":
				// send answer to game and proceed to new questoin

			}
		}
	}(conn)

	logger.L().Info("Successfuly started websocket.", zap.Any("gorutoines", runtime.NumGoroutine()))

	return nil
}

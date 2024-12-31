package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/claim"
	"mdhesari/kian-quiz-golang-game/pkg/timestamp"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())
	if err != nil {
		logger.L().Error("Could not upgrade http to websocket.", zap.Error(err))

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	go func(conn *websocket.Conn) {
		defer conn.Close()

		ticker := time.NewTicker(60 * time.Second)
		defer func() {
			ticker.Stop()
			conn.Close()
		}()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
					logger.L().Error("Error sending ping.", zap.Error(err))

					return
				}

				logger.L().Info("ping")
			}
		}
	}(conn)

	go func(conn *websocket.Conn) {
		defer conn.Close()

		conn.SetPongHandler(func(msg string) error {
			if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
				return err
			}

			userId := claims.UserID
			_, err := h.presenceSrv.Upsert(c.Request().Context(), param.PresenceUpsertRequest{
				UserId:    userId,
				Timestamp: timestamp.Now(),
			})
			if err != nil {
				logger.L().Error("Error upserting presence after pong.")

				return err
			}

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
		}
	}(conn)

	logger.L().Info("Successfuly started websocket.", zap.Any("gorutoines", runtime.NumGoroutine()))

	return nil
}

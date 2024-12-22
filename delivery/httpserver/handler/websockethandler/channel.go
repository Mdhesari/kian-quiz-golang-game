package websockethandler

import (
	"context"
	"mdhesari/kian-quiz-golang-game/logger"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *Handler) Channel(c echo.Context) error {
	conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response().Writer)
	if err != nil {
		logger.L().Error("Could not upgrade http to websocket.", zap.Error(err))

		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer conn.Close()

	channel := c.Param("channel")
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	pubsub := h.redisAdap.Subscribe(c.Request().Context(), channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	go func() {
		defer conn.Close()
		defer pubsub.Close()

		for {
			msg, opCode, err := wsutil.ReadClientData(conn)
			if err != nil {
				logger.L().Error("Could not read client data.", zap.Error(err), zap.String("opcode", string(opCode)))
				cancel()

				return
			}

			h.redisAdap.Publish(ctx, channel, string(msg))
		}
	}()

	for msg := range ch {
		err := wsutil.WriteServerMessage(conn, ws.OpText, []byte(msg.Payload))
		if err != nil {
			logger.L().Error("Could not websocket server message.", zap.Error(err), zap.Any("payload", msg.PayloadSlice))

			cancel()

			break
		}
	}

	return nil
}

package websockethandler

import (
	"context"
	"mdhesari/kian-quiz-golang-game/logger"
	"net/http"
	"time"

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

	if err := conn.SetReadDeadline(time.Now().Add(h.cfg.ReadTimeout)); err != nil {
		logger.L().Error("Coud not set read deadline for websocket.", zap.Error(err))

		return err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(h.cfg.WriteTimeout)); err != nil {
		logger.L().Error("Coud not set write deadline for websocket.", zap.Error(err))

		return err
	}

	channel := c.Param("channel")
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	pubsub := h.redisAdap.Subscribe(c.Request().Context(), channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	go func() {
		for {
			msg, opCode, err := wsutil.ReadClientData(conn)
			if err != nil {
				logger.L().Error("Could not read client data.", zap.Error(err), zap.String("opcode", string(opCode)))
				cancel()

				return
			}

			if opCode == ws.OpPing {
				if err := wsutil.WriteClientMessage(conn, ws.OpPong, nil); err != nil {
					logger.L().Error("Could not write client message pong.")
					cancel()

					return
				}

				continue
			}

			h.redisAdap.Publish(ctx, channel, string(msg))
		}
	}()

	ticker := time.NewTicker(h.cfg.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.L().Info("Websocket connection closed due to context cancellation.")

			return nil
		case msg := <-ch:
			err := wsutil.WriteServerMessage(conn, ws.OpText, []byte(msg.Payload))
			if err != nil {
				logger.L().Error("Could not websocket server message.", zap.Error(err), zap.Any("payload", msg.PayloadSlice))

				cancel()

				return err
			}
		case <-ticker.C:
			if err := wsutil.WriteClientMessage(conn, ws.OpPing, nil); err != nil {
				cancel()

				return err
			}
		}
	}
}

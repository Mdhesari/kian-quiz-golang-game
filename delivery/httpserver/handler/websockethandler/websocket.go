package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/logger"
	"net/http"
	"runtime"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h Handler) Websocket(c echo.Context) error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), c.Response().Header())
	if err != nil {
		logger.L().Error("Could not upgrade http to websocket.", zap.Error(err))

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	cli := h.websocketAdapt.NewClient(conn)
	h.websocketAdapt.RegisterClient(cli)

	go cli.ReadPump()
	go cli.WritePump()

	logger.L().Info("Successfuly started websocket.", zap.Any("gorutoines", runtime.NumGoroutine()))

	return nil
}

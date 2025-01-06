package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/claim"
	"mdhesari/kian-quiz-golang-game/websockethub"
	"net/http"
	"runtime"

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
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		logger.L().Error("Could not upgrade http to websocket.", zap.Error(err))

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	cli := websockethub.NewClient(conn, h.hub, claims.UserID)
	cli.Start()

	h.hub.RegisterClient(&cli)

	logger.L().Info("Successfuly started websocket.", zap.Any("gorutoines", runtime.NumGoroutine()))

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func producer(remoteAddr string, channel chan<- string) {
	for {
		channel <- remoteAddr
		channel <- "pong"
		time.Sleep(5 * time.Second)
	}
}

func main() {
	http.ListenAndServe(":9000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
			panic(err)
		}

		defer conn.Close()

		done := make(chan bool)
		go readMessage(conn, done)

		channel := make(chan string)
		go producer(r.RemoteAddr, channel)
		go writeMessage(conn, channel)

		<-done

		fmt.Println("cli disconnected.")
	}))
}

func readMessage(conn net.Conn, done chan<- bool) {
	for {
		msg, opCode, err := wsutil.ReadClientData(conn)
		if err != nil {
			// TODO - update metrics
			if err != io.EOF {
				log.Fatalf("could not read message {%v}.\n", err)
			}

			done <- true

			return
		}
		if opCode == ws.OpClose {

			return
		}
		fmt.Println("op code", opCode, msg)

		var notif entity.Notification
		err = json.Unmarshal(msg, &notif)
		if err != nil {
			panic(err)
		}

		fmt.Println("opCode", opCode)
		fmt.Println("notif", notif)
	}
}

func writeMessage(conn net.Conn, channel <-chan string) {
	for data := range channel {
		err := wsutil.WriteServerMessage(conn, ws.OpText, []byte(data))
		if err != nil {

			fmt.Println("Client disconnected")
		}
	}
}

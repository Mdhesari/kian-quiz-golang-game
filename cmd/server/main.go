package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mdhesari/kian-quiz-golang-game/repository/mongorepo"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Everything is fine!"))
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Write([]byte(`{"errors": ["Not supported method."]}`))

			return
		}

		repo := mongorepo.New("mongodb://michael:secret@db:27017")

		usersrv := userservice.New(repo)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Request body read error: %w", err)

			w.Write([]byte("{}"))

			return
		}

		uf := userservice.UserForm{}
		json.Unmarshal(body, &uf)

		user, err := usersrv.Register(uf)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"errors": ["%s"]}`, err.Error())))

			return
		}

		w.Write([]byte(fmt.Sprintf(`{"message" : "User %s is created."}`, user.Name)))
	})

	server := http.Server{Addr: ":2001", Handler: mux}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

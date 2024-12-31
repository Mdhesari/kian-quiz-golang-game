package websocketadapter

import "github.com/redis/go-redis/v9"

type Hub struct {
	channel <-chan *redis.Message
}

func NewHub() *Hub {
	return &Hub{
		//
	}
}

func (h *Hub) Run() {
	for {
		select {
		case msg := <- h.channel:
			
		}
	}
}

package entity

type WebsocketMsg struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

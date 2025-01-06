package entity

type WebsocketMsg struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

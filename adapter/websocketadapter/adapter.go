package websocketadapter

type Config struct {
	ClientCfg ClientConfig `koanf:"client"`
}

type Adapter struct {
	cfg Config
	hub *Hub
}

func New(cfg Config) *Adapter {
	return &Adapter{
		cfg: cfg,
		hub: NewHub(),
	}
}

func (a *Adapter) RegisterClient(client *Client) {
	a.hub.register <- client
}

func (a *Adapter) Run() {
	a.hub.Run()
}

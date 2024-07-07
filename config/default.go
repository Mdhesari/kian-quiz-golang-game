package config

var defaultConfig = map[string]interface{}{
	"redis": map[string]interface{}{
		"host":     "127.0.0.1",
		"port":     6379,
		"password": "",
		"database": 0,
		"username": 0,
	},
}

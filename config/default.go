package config

import "time"

var defaultConfig = map[string]interface{}{
	"redis": map[string]interface{}{
		"host":     "127.0.0.1",
		"port":     6379,
		"password": "",
		"database": 0,
		"username": 0,
	},
	"application": map[string]interface{}{
		"graceful_shutdown_timeout": 5 * time.Second,
	},
}

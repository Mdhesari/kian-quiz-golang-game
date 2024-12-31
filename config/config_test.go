package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {
	cfg := Load("./../config.test.yml")

	tests := []struct {
		expected interface{}
		actual   interface{}
	}{
		{
			27017,
			cfg.Database.MongoDB.Port,
		},
		{
			"mongodb",
			cfg.Database.MongoDB.Host,
		},
		{
			"db",
			cfg.Database.MongoDB.DBName,
		},
		{
			"michael",
			cfg.Database.MongoDB.Username,
		},
		{
			"file://repository/mongorepo/seeders",
			cfg.Database.Seeders,
		},
		{
			"file://repository/mongorepo/migrations",
			cfg.Database.Migrations,
		},
		{
			5,
			cfg.Database.MongoDB.DurationSeconds,
		},
		{
			"redis",
			cfg.Redis.Host,
		},
		{
			6,
			cfg.Application.Question.QuestionsCount,
		},
		{
			80,
			cfg.Server.HttpServer.Port,
		},
		{
			[]byte("secret"),
			cfg.Auth.Secret,
		},
		{
			time.Duration(10080),
			cfg.Auth.ExpiresInMinutes,
		},
		{
			"presence",
			cfg.Presence.Prefix,
		},
		{
			time.Duration(4 * time.Hour),
			cfg.Presence.Expiration,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.actual)
	}
}

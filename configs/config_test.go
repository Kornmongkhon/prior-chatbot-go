package configs

import (
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

func TestGetConfig(t *testing.T) {
	err := Init(".")
	if err != nil {
		panic(err)
	}
	cfg := GetConfig()

	assert.Equal(t, "prior-chatbot", cfg.Server.Name)
	assert.Equal(t, "prior_chatbot", cfg.DB.Database)
	assert.Equal(t, "Asia/Bangkok", cfg.Server.TimeZone)
	assert.Equal(t, time.Duration(900000000000), cfg.JWT.ExpirationAccessToken)
	assert.Equal(t, "prior-chatbot-secret", cfg.JWT.Secret)

}

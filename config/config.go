package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	webdock "github.com/webdock-io/go-sdk"
	"github.com/zolamk/terraform-provider-webdock/api"
)

type Config struct {
	Token            string
	TerraformVersion string
	ServerUpPort     int
	RetryLimit       int
}

type Counter struct {
	mu sync.Mutex
	x  int64
}

func (c *Counter) Inc() {
	c.mu.Lock()
	c.x++
	c.mu.Unlock()
}

func (c *Counter) Value() (x int64) {
	c.mu.Lock()
	x = c.x
	c.mu.Unlock()
	return
}

type CombinedConfig struct {
	api.ClientInterface
	Logger              *slog.Logger
	CreatedServersCount Counter
	CreateUsersCount    Counter
	ServerUpPort        int
	RetryLimit          int
}

func NewCombinedConfig(config *Config, client api.ClientInterface) *CombinedConfig {
	return &CombinedConfig{
		client,
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Counter{},
		Counter{},
		config.ServerUpPort,
		config.RetryLimit,
	}
}

func (c *Config) Client() (*CombinedConfig, diag.Diagnostics) {
	client := webdock.New(webdock.WebdockOptions{TOKEN: c.Token})

	return &CombinedConfig{
		&client,
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Counter{},
		Counter{},
		c.ServerUpPort,
		c.RetryLimit,
	}, nil
}

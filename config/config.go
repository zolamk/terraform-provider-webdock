package config

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/zolamk/terraform-provider-webdock/api"
)

type Config struct {
	Token            string
	APIEndpoint      string
	TerraformVersion string
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
}

func NewCombinedConfig(config *Config, client api.ClientInterface) *CombinedConfig {
	return &CombinedConfig{
		client,
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Counter{},
	}
}

func setAuthorization(c *Config) api.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+c.Token)
		return nil
	}
}

func (c *Config) Client() (*CombinedConfig, diag.Diagnostics) {
	webdockClient, err := api.NewClient(c.APIEndpoint+"/v1", api.WithRequestEditorFn(setAuthorization(c)))
	if err != nil {
		return nil, diag.Errorf("error creating api client: %v", err)
	}

	return &CombinedConfig{
		webdockClient,
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Counter{},
	}, nil
}

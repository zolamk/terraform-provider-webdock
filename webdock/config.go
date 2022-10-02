package webdock

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/api"
)

type Config struct {
	Token            string
	APIEndpoint      string
	TerraformVersion string
}

type CombinedConfig struct {
	client *api.Client
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
		client: webdockClient,
	}, nil
}

// waitForAction waits for the action to finish using the resource.StateChangeConf.
func waitForAction(client *api.Client, callbackID string) error {
	var (
		pending   = "waiting"
		working   = "working"
		target    = "finished"
		refreshfn = func() (result interface{}, state string, err error) {

			opts := &api.GetEventsParams{}

			events, err := client.GetEvents(context.Background(), opts)

			if err != nil {
				return nil, "", err
			}

			if events == nil || len(*events) == 0 {
				return nil, "", errors.New("error getting event state: response body empty")
			}

			event := (*events)[0]

			return event, string(event.Status), nil

		}
	)
	_, err := (&resource.StateChangeConf{
		Pending:    []string{pending, working},
		Refresh:    refreshfn,
		Target:     []string{target},
		Delay:      10 * time.Second,
		Timeout:    60 * time.Minute,
		MinTimeout: 3 * time.Second,
	}).WaitForState()

	return err
}
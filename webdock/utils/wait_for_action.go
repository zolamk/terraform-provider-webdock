package utils

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zolamk/terraform-provider-webdock/api"
)

func WaitForAction(ctx context.Context, client api.ClientInterface, callbackID string) error {
	var (
		pending   = "waiting"
		working   = "working"
		target    = "finished"
		refreshfn = func() (result interface{}, state string, err error) {
			opts := api.GetEventsParams{
				CallbackId: callbackID,
			}

			events, err := client.GetEvents(ctx, opts)
			if err != nil {
				return nil, "", err
			}

			if len(events) == 0 {
				return nil, "", errors.New("error getting event state: response body empty")
			}

			event := (events)[0]

			return event, event.Status, nil
		}
	)
	_, err := (&resource.StateChangeConf{
		Pending:    []string{pending, working},
		Refresh:    refreshfn,
		Target:     []string{target},
		Delay:      10 * time.Second,
		Timeout:    60 * time.Minute,
		MinTimeout: 3 * time.Second,
	}).WaitForStateContext(ctx)

	return err
}

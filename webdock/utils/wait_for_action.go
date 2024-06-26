package utils

import (
	"context"
	"errors"
	"fmt"
	"net"
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
		Timeout:    10 * time.Minute,
		MinTimeout: 3 * time.Second,
	}).WaitForStateContext(ctx)

	return err
}

// WaitForServerToBeUp makes sure besides of getting finished status that the server is actually reachable on port 22
func WaitForServerToBeUP(ctx context.Context, client api.ClientInterface, callbackID string, ip string, port int) error {
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

			if event.Status == target {
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Minute)
				if err != nil {
					return event, working, nil
				}

				defer conn.Close()
			}

			return event, event.Status, nil
		}
	)
	_, err := (&resource.StateChangeConf{
		Pending:    []string{pending, working},
		Refresh:    refreshfn,
		Target:     []string{target},
		Delay:      10 * time.Second,
		Timeout:    10 * time.Minute,
		MinTimeout: 3 * time.Second,
	}).WaitForStateContext(ctx)

	return err
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// Event log status
type EventLogStatus string

// Defines values for EventLog Status.
const (
	EventLogDTOStatusError    EventLogStatus = "error"
	EventLogDTOStatusFinished EventLogStatus = "finished"
	EventLogDTOStatusWaiting  EventLogStatus = "waiting"
	EventLogDTOStatusWorking  EventLogStatus = "working"
)

// GetEventsParams defines parameters for GetEvents.
type GetEventsParams struct {
	// Callback ID
	CallbackId string `json:"callbackId,omitempty"`

	// Event Type
	EventType string `json:"eventType,omitempty"`

	// Page
	Page int64 `json:"page,omitempty"`

	// Events per page
	PerPage int64 `json:"per_page,omitempty"`
}

// Event log model
type EventLog struct {
	// Just a plain text description of the action. Same text as you see in the Event Log in the Webdock Dashboard.
	Action string `json:"action,omitempty"`

	// Action Data. A more static/parseable string representation of the action.
	ActionData string `json:"actionData,omitempty"`

	// Callback ID
	CallbackId string `json:"callbackId,omitempty"`

	// End Time of the event
	EndTime string `json:"endTime"`

	// Event Type
	EventType string `json:"eventType,omitempty"`

	// Event log ID
	Id json.Number `json:"id,omitempty"`

	// Any &quot;Message&quot; or return data from the action once finished executing.
	Message string `json:"message,omitempty"`

	// Server Slug
	ServerSlug string `json:"serverSlug,omitempty"`

	// Start Time of the event
	StartTime string `json:"startTime,omitempty"`

	// Status
	Status EventLogStatus `json:"status,omitempty"`
}

type Events []EventLog

func (c *Client) GetEvents(ctx context.Context, params *GetEventsParams) (Events, error) {
	serverURL, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	serverURL.Path += "events"

	queryValues, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	serverURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error getting events: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting events: %w", err)
	}

	defer resp.Body.Close()

	if errorStatus(resp.StatusCode) {
		apiError := APIError{}

		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error decoding get events error response body: %w", err)
		}

		return nil, fmt.Errorf("error getting events: %w", apiError)
	}

	events := Events{}

	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("error decoding get events response body: %w", err)
	}

	return events, nil
}

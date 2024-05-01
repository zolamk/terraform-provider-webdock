package resource

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/api"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
	"github.com/zolamk/terraform-provider-webdock/webdock/utils"
)

var (
	tooManyServersMessage = "You are creating too many servers in too short of a timespan. Please wait a while and try again a bit later."
)

func Server() *schema.Resource {
	return &schema.Resource{
		CreateContext: createServer,
		ReadContext:   readServer,
		UpdateContext: updateServer,
		DeleteContext: deleteServer,
		SchemaVersion: 0,
		Schema:        schemas.Server(),
	}
}

func createServer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	retryCount := 3
	retriedTimes := 0

	client := meta.(*config.CombinedConfig)

	delay := time.Duration(client.CreatedServersCount.Value()*10) * time.Second

	client.CreatedServersCount.Inc()

	client.Logger.With("delay", delay).Info("sleeping to avoid concurrency issues")

	time.Sleep(delay)

	opts := api.CreateServerRequestBody{
		Name:           d.Get("name").(string),
		LocationId:     d.Get("location_id").(string),
		ProfileSlug:    d.Get("profile_slug").(string),
		ImageSlug:      d.Get("image_slug").(string),
		Virtualization: d.Get("virtualization").(string),
	}

	if attr, ok := d.GetOk("slug"); ok {
		opts.Slug = attr.(string)
	}

createServer:
	server, err := client.CreateServer(ctx, opts)
	if err != nil {
		apiErr := &api.APIError{}
		if (errors.As(err, apiErr) || errors.As(err, &apiErr)) && apiErr.Message == tooManyServersMessage && retriedTimes < retryCount {
			retriedTimes++

			delay := time.Duration(retriedTimes) * time.Minute

			client.Logger.With("delay", delay).Info("got too many servers error, will retry after a while")

			time.Sleep(delay)

			goto createServer
		}

		return diag.FromErr(err)
	}

	d.SetId(server.Slug)

	if server.CallbackID != "" {
		err = utils.WaitForServerToBeUP(ctx, client, server.CallbackID, server.Ipv4, client.ServerUpPort)
		if err != nil {
			return diag.Errorf("server (%s) create event (%s) errored: %v", d.Id(), server.CallbackID, err)
		}
	} else {
		return diag.Errorf("unable to find server (%s) create event", d.Id())
	}

	if err := setServerAttributes(d, server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readServer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	server, err := client.GetServerBySlug(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, api.ErrServerNotFound) {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting server: %v", err)
	}

	if err = setServerAttributes(d, server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updateServer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	if d.HasChange("profile_slug") {
		_, newProfileSlug := d.GetChange("profile_slug")

		opts := api.ResizeServerRequestBody{
			ProfileSlug: newProfileSlug.(string),
		}

		_, err := client.ResizeDryRun(ctx, d.Id(), opts)
		if err != nil {
			return diag.FromErr(err)
		}

		callbackID, err := client.ResizeServer(ctx, d.Id(), opts)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = utils.WaitForAction(ctx, client, callbackID); err != nil {
			return diag.Errorf("server (%s) profile change event (%s) errorred: %s", d.Id(), callbackID, err)
		}
	}

	if d.HasChange("image_slug") {
		_, newImageSlug := d.GetChange("image_slug")

		opts := api.ReinstallServerRequestBody{
			ImageSlug: newImageSlug.(string),
		}

		callbackID, err := client.ReinstallServer(ctx, d.Id(), opts)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = utils.WaitForAction(ctx, client, callbackID); err != nil {
			return diag.Errorf("server (%s) reinstall event (%s) errorred: %s", d.Id(), callbackID, err)
		}
	}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")

		opts := api.PatchServerRequestBody{
			Name: newName.(string),
		}

		if _, err := client.PatchServer(ctx, d.Id(), opts); err != nil {
			return diag.FromErr(err)
		}
	}

	return readServer(ctx, d, meta)
}

func deleteServer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	callbackID, err := client.DeleteServer(context.Background(), d.Id())

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return nil
		}

		return diag.FromErr(err)
	}

	if err = utils.WaitForAction(ctx, client, callbackID); err != nil {
		diag.Errorf("server (%s) delete event (%s) errorred: %s", d.Id(), callbackID, err)
	}

	d.SetId("")

	return nil
}

func setServerAttributes(d *schema.ResourceData, server *api.Server) error {
	if err := d.Set("name", server.Name); err != nil {
		return err
	}

	if err := d.Set("slug", server.Slug); err != nil {
		return err
	}

	if err := d.Set("location_id", server.Location); err != nil {
		return err
	}

	if err := d.Set("profile_slug", server.Profile); err != nil {
		return err
	}

	if err := d.Set("image_slug", server.Image); err != nil {
		return err
	}

	if err := d.Set("created_at", server.Date); err != nil {
		return err
	}

	if err := d.Set("ipv4", server.Ipv4); err != nil {
		return err
	}

	if err := d.Set("ipv6", server.Ipv6); err != nil {
		return err
	}

	if err := d.Set("status", server.Status); err != nil {
		return err
	}

	if err := d.Set("webserver", server.WebServer); err != nil {
		return err
	}

	if err := d.Set("aliases", server.Aliases); err != nil {
		return err
	}

	if err := d.Set("snapshot_runtime", server.SnapshotRunTime); err != nil {
		return err
	}

	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": server.Ipv4,
	})

	return nil
}

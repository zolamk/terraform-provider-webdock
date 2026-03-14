package resource

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	webdock "github.com/webdock-io/go-sdk"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Hour * 3),
		},
	}
}

func createServer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	currentAttempt := 0

	initialInterval := 1 * time.Minute

	client := meta.(*config.CombinedConfig)

	delay := time.Duration(client.CreatedServersCount.Value()*10) * time.Second

	client.CreatedServersCount.Inc()

	client.Logger.With("delay", delay).Info("sleeping to avoid concurrency issues")

	time.Sleep(delay)

	opts := webdock.CreateServerFromImageOptions{
		Name:        d.Get("name").(string),
		LocationId:  d.Get("location_id").(string),
		ProfileSlug: d.Get("profile_slug").(string),
		ImageSlug:   d.Get("image_slug").(string),
	}

createServer:
	createdServer, err := client.CreateServerFromImage(opts)
	if err != nil {
		if strings.Contains(err.Error(), tooManyServersMessage) && currentAttempt < client.RetryLimit {
			currentAttempt++

			delay := initialInterval * time.Duration(math.Pow(2, float64(currentAttempt-1)))

			client.Logger.With("delay", delay).Info("got too many servers error, will retry after a while")

			time.Sleep(delay)

			goto createServer
		}

		return diag.FromErr(err)
	}

	server := createdServer.Server
	d.SetId(server.Slug)

	err = utils.WaitForServerToBeUP(ctx, client, createdServer.CallbackID, server.IPv4, client.ServerUpPort)
	if err != nil {
		return diag.Errorf("server (%s) create event (%s) errored: %v", d.Id(), createdServer.CallbackID, err)
	}

	if err := setServerAttributes(d, &server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func readServer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	server, err := client.GetServerBySlug(webdock.GetServerBySlugOptions{Slug: d.Id()})
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") || strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting server: %v", err)
	}

	if err = setServerAttributes(d, &server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updateServer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	if d.HasChange("profile_slug") {
		_, newProfileSlug := d.GetChange("profile_slug")

		opts := webdock.ResizeServersOptions{
			Slug:        d.Id(),
			ProfileSlug: newProfileSlug.(string),
		}

		dryOpts := webdock.DryRunResizeServerOptions{
			ServerSlug:  d.Id(),
			ProfileSlug: newProfileSlug.(string),
		}

		_, err := client.DryRunResizeServer(dryOpts)
		if err != nil {
			return diag.FromErr(err)
		}

		callbackID, err := client.ResizeServer(opts)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = utils.WaitForAction(ctx, client, callbackID); err != nil {
			return diag.Errorf("server (%s) profile change event (%s) errorred: %s", d.Id(), callbackID, err)
		}
	}

	if d.HasChange("image_slug") {
		_, newImageSlug := d.GetChange("image_slug")

		opts := webdock.ReinstallServerOptions{
			Slug:      d.Id(),
			ImageSlug: newImageSlug.(string),
		}

		callbackID, err := client.ReinstallServer(opts)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = utils.WaitForAction(ctx, client, callbackID); err != nil {
			return diag.Errorf("server (%s) reinstall event (%s) errorred: %s", d.Id(), callbackID, err)
		}
	}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")

		opts := webdock.UpdateServerOptions{
			ServerSlug: d.Id(),
			Name:       newName.(string),
		}

		if _, err := client.UpdateServer(opts); err != nil {
			return diag.FromErr(err)
		}
	}

	return readServer(ctx, d, meta)
}

func deleteServer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	err := client.DeleteServerBySlug(webdock.DeleteServerBySlugOptions{Slug: d.Id()})

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") || strings.Contains(err.Error(), "404") {
			return nil
		}

		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setServerAttributes(d *schema.ResourceData, server *webdock.Server) error {
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

	if err := d.Set("ipv4", server.IPv4); err != nil {
		return err
	}

	if err := d.Set("ipv6", server.IPv6); err != nil {
		return err
	}

	if err := d.Set("status", string(server.Status)); err != nil {
		return err
	}

	if err := d.Set("webserver", string(server.WebServer)); err != nil {
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
		"host": server.IPv4,
	})

	return nil
}

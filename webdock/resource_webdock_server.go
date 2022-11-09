package webdock

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/api"
)

func resourceWebdockServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebdockServerCreate,
		ReadContext:   resourceWebdockServerRead,
		UpdateContext: resourceWebdockServerUpdate,
		DeleteContext: resourceWebdockServerDelete,
		SchemaVersion: 0,
		Schema:        serverSchema(),
	}
}

func resourceWebdockServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

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

	server, err := client.CreateServer(context.Background(), opts)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(server.Slug)

	if server.CallbackID != "" {
		err = waitForAction(client, server.CallbackID)
		if err != nil {
			return diag.Errorf("server (%s) create event (%s) errorred: %s", d.Id(), server.CallbackID, err)
		}
	} else {
		return diag.Errorf("unable to find server (%s) create event.", d.Id())
	}

	if err := setServerAttributes(d, server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceWebdockServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	server, err := client.GetServerBySlug(context.Background(), d.Id())
	if err != nil {
		// check if the server no longer exists
		if strings.Contains(err.Error(), "Not Found") {
			log.Printf("[WARN] Webdock server (%s) not found", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving server: %s", err)
	}

	if err = setServerAttributes(d, server); err != nil {
		return diag.FromErr(err)
	}

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

	if err := d.Set("description", server.Description); err != nil {
		return err
	}

	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": server.Ipv4,
	})

	return nil
}

func resourceWebdockServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	if d.HasChange("profile_slug") {
		_, newProfileSlug := d.GetChange("profile_slug")

		opts := api.ResizeServerRequestBody{
			ProfileSlug: newProfileSlug.(string),
		}

		_, err := client.ResizeDryRun(context.Background(), d.Id(), opts)

		if err != nil {
			return diag.FromErr(err)
		}

		callbackID, err := client.ResizeServer(context.Background(), d.Id(), opts)

		if err != nil {
			return diag.FromErr(err)
		}

		if err = waitForAction(client, callbackID); err != nil {
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

		if err = waitForAction(client, callbackID); err != nil {
			return diag.Errorf("server (%s) reinstall event (%s) errorred: %s", d.Id(), callbackID, err)
		}
	}

	opts := api.PatchServerRequestBody{
		Name: d.Get("name").(string),
	}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")

		opts.Name = newName.(string)
	}

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")

		opts.Description = newDescription.(string)
	}

	if d.HasChange("notes") {
		_, newNotes := d.GetChange("notes")

		opts.Notes = newNotes.(string)
	}

	if d.HasChange("next_action_date") {
		_, newNextActionDate := d.GetChange("next_action_date")

		opts.NextActionDate = newNextActionDate.(string)
	}

	if _, err := client.PatchServer(context.Background(), d.Id(), opts); err != nil {
		return diag.FromErr(err)
	}

	return resourceWebdockServerRead(ctx, d, meta)
}

func resourceWebdockServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	callbackID, err := client.DeleteServer(context.Background(), d.Id())

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return nil
		}

		return diag.FromErr(err)
	}

	if err = waitForAction(client, callbackID); err != nil {
		diag.Errorf("server (%s) delete event (%s) errorred: %s", d.Id(), callbackID, err)
	}

	d.SetId("")

	return nil
}

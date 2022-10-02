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

	// Build up our creation options
	opts := api.CreateServerJSONRequestBody{
		Name:        d.Get("name").(string),
		LocationId:  d.Get("location_id").(string),
		ProfileSlug: d.Get("profile_slug").(string),
		ImageSlug:   d.Get("image_slug").(string),
	}

	if attr, ok := d.GetOk("slug"); ok {
		opts.Slug = attr.(*string)
	}

	log.Printf("[DEBUG] Server create configuration: %#v", opts)

	server, callbackID, err := client.CreateServer(context.Background(), opts)

	if err != nil {
		return diag.Errorf("Error creating server: %s", err)
	}

	d.SetId(server.Slug)

	// Wait for server create action to successfully finish.
	if *callbackID != "" {
		err = waitForAction(client, *callbackID)
		if err != nil {
			return diag.Errorf("Server (%s) create event (%d) errorred: %s", d.Id(), callbackID, err)
		}
	} else {
		return diag.Errorf("Unable to find server (%s) create event.", d.Id())
	}

	setServerAttributes(d, server)

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
	d.Set("name", server.Name)

	d.Set("slug", server.Slug)

	d.Set("location_id", server.Location)

	d.Set("profile_slug", server.Profile)

	d.Set("image_slug", server.Image)

	d.Set("created_at", server.Date)

	d.Set("ipv4", server.Ipv4)

	d.Set("ipv6", server.Ipv6)

	d.Set("status", server.Status)

	d.Set("webserver", server.WebServer)

	d.Set("aliases", server.Aliases)

	d.Set("snapshot_runtime", server.SnapshotRunTime)

	d.Set("description", server.Description)

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

		opts := api.ResizeServerModelDTO{
			ProfileSlug: newProfileSlug.(string),
		}

		_, err := client.ResizeDryRun(context.Background(), d.Id(), opts)

		if err != nil {
			return diag.Errorf("Error changing server profile: %v", err)
		}

		callbackID, err := client.ResizeServer(context.Background(), d.Id(), opts)

		if err != nil {
			return diag.Errorf("Error changing server profile: %v", err)
		}

		if err = waitForAction(client, *callbackID); err != nil {
			return diag.Errorf("Server (%s) profile change event (%d) errorred: %s", d.Id(), callbackID, err)
		}
	}

	if d.HasChange("image_slug") {
		_, newImageSlug := d.GetChange("image_slug")

		opts := api.ReinstallServerModelDTO{
			ImageSlug: newImageSlug.(string),
		}

		callbackID, err := client.ReinstallServer(ctx, d.Id(), opts)

		if err != nil {
			return diag.Errorf("Error reinstalling server: %v", err)
		}

		if err = waitForAction(client, *callbackID); err != nil {
			return diag.Errorf("Server (%s) reinstall event (%d) errorred: %s", d.Id(), callbackID, err)
		}
	}

	opts := api.PatchServerJSONRequestBody{
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

		opts.NextActionDate = newNextActionDate.(*string)
	}

	if _, err := client.PatchServer(context.Background(), d.Id(), opts); err != nil {
		return diag.Errorf("Error updating server (%s): %s", d.Id(), err)
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

		return diag.Errorf("Error deleting server (%s): %s", d.Id(), err)
	}

	if err = waitForAction(client, *callbackID); err != nil {
		return diag.Errorf("Error deleting server (%s): %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

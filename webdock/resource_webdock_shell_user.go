package webdock

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/api"
)

func resourceWebdockShellUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebdockShellUserCreate,
		UpdateContext: resourceWebdockShellUserUpdate,
		DeleteContext: resourceWebdockShellUserDelete,
		ReadContext:   resourceWebdockShellUserRead,
		SchemaVersion: 0,
		Schema:        shellUserSchema(),
	}
}

func resourceWebdockShellUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	shellUsers, err := client.GetShellUsers(ctx, d.Get("server_slug").(string))
	if err != nil {
		return diag.Errorf("error getting shell users: %v", err)
	}

	shellUser := findShellUserByID(d.Id(), shellUsers)

	if shellUser == nil {
		return diag.Errorf("error getting public key: 404 Not Found")
	}

	if err = setShellUserAttributes(d, shellUser); err != nil {
		return diag.Errorf("error setting public key: %v", err)
	}

	return nil
}

func resourceWebdockShellUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	var publicKeys []int64

	for _, key := range d.Get("public_keys").([]interface{}) {
		publicKeys = append(publicKeys, key.(int64))
	}

	createShellUserBody := api.CreateShellUserRequestBody{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		Group:      d.Get("group").(string),
		Shell:      d.Get("shell").(string),
		PublicKeys: publicKeys,
	}

	shellUser, err := client.CreateShellUser(ctx, d.Get("server_slug").(string), createShellUserBody)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := waitForAction(ctx, client, shellUser.CallbackID); err != nil {
		return diag.Errorf("error creating shell user: %s", err)
	}

	if err = setShellUserAttributes(d, shellUser); err != nil {
		return diag.Errorf("error setting shell user: %s", err)
	}

	return nil
}

func resourceWebdockShellUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("error converting id to number: %v", err)
	}

	shellUser, err := client.UpdateShellUserPublicKeys(ctx, d.Get("server_slug").(string), id, d.Get("public_keys").([]int64))
	if err != nil {
		return diag.Errorf("error updating shell user: %v", err)
	}

	if err := waitForAction(ctx, client, shellUser.CallbackID); err != nil {
		return diag.Errorf("error updating shell user: %v", err)
	}

	if err := d.Set("public_keys", shellUser.PublicKeys); err != nil {
		return diag.Errorf("error setting public keys: %v", err)
	}

	return nil
}

func resourceWebdockShellUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("error converting id to number: %v", err)
	}

	callbackID, err := client.DeleteShellUser(ctx, d.Get("server_slug").(string), id)
	if err != nil {
		return diag.Errorf("error deleting shell user: %v", err)
	}

	if err = waitForAction(ctx, client, callbackID); err != nil {
		return diag.Errorf("Error deleting shell user (%s): %v", d.Id(), err)
	}

	return nil
}

func setShellUserAttributes(d *schema.ResourceData, shellUser *api.ShellUser) error {
	d.SetId(shellUser.ID.String())

	if err := d.Set("username", shellUser.Username); err != nil {
		return err
	}

	if err := d.Set("group", shellUser.Group); err != nil {
		return err
	}

	if err := d.Set("shell", shellUser.Shell); err != nil {
		return err
	}

	if err := d.Set("created_at", shellUser.Created); err != nil {
		return err
	}

	return nil
}

func findShellUserByID(id string, shellUsers api.ShellUsers) *api.ShellUser {
	if shellUsers == nil {
		return nil
	}

	for _, shellUser := range shellUsers {
		if shellUser.ID.String() == id {
			return &shellUser
		}
	}

	return nil
}

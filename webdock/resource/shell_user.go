package resource

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	webdock "github.com/webdock-io/go-sdk"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
	"github.com/zolamk/terraform-provider-webdock/webdock/utils"
)

func ShellUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: createShellUser,
		UpdateContext: updateShellUser,
		DeleteContext: deleteShellUser,
		ReadContext:   readShellUser,
		SchemaVersion: 0,
		Schema:        schemas.ShellUser(),
	}
}

func createShellUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	var publicKeys []int64

	for _, key := range d.Get("public_keys").([]interface{}) {
		publicKeys = append(publicKeys, int64(key.(int)))
	}

	delay := time.Duration(client.CreateUsersCount.Value()*10) * time.Second

	client.CreateUsersCount.Inc()

	client.Logger.With("delay", delay).Info("sleeping to avoid concurrency issues")

	time.Sleep(delay)

	opts := webdock.CreateServerShellUserOptions{
		ServerSlug: d.Get("server_slug").(string),
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		Group:      d.Get("group").(string),
		Shell:      d.Get("shell").(string),
		PublicKeys: publicKeys,
	}

	createdShellUser, err := client.CreateServerShellUser(opts)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := utils.WaitForAction(ctx, client, createdShellUser.CallbackID); err != nil {
		return diag.Errorf("error creating shell user: %s", err)
	}

	shellUser := createdShellUser.ShellUser

	if err = setShellUserAttributes(d, &shellUser); err != nil {
		return diag.Errorf("error setting shell user: %s", err)
	}

	return nil
}

func readShellUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	opts := webdock.ListServerShellUserOptions{
		ServerSlug: d.Get("server_slug").(string),
	}

	shellUsers, err := client.ListServerShellUser(opts)
	if err != nil {
		return diag.Errorf("error getting shell users: %v", err)
	}

	shellUser := findShellUserByID(d.Id(), shellUsers)

	if shellUser == nil {
		return diag.Errorf("error getting shell user: 404 Not Found")
	}

	if err = setShellUserAttributes(d, shellUser); err != nil {
		return diag.Errorf("error setting shell user: %v", err)
	}

	return nil
}

func updateShellUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("error converting id to number: %v", err)
	}

	var publicKeys []int64
	for _, key := range d.Get("public_keys").([]interface{}) {
		publicKeys = append(publicKeys, int64(key.(int)))
	}

	opts := webdock.UpdateServerShellUserOptions{
		ServerSlug:  d.Get("server_slug").(string),
		ShellUserId: id,
		PublicKeys:  publicKeys,
	}

	shellUser, err := client.UpdateServerShellUser(opts)
	if err != nil {
		return diag.Errorf("error updating shell user: %v", err)
	}

	if err := d.Set("public_keys", shellUser.PublicKeys); err != nil {
		return diag.Errorf("error setting public keys: %v", err)
	}

	return nil
}

func deleteShellUser(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("error converting id to number: %v", err)
	}

	opts := webdock.DeleteShellUserOptions{
		ServerSlug:  d.Get("server_slug").(string),
		ShellUserId: id,
	}

	deleteRes, err := client.DeleteShellUser(opts)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") || strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.Errorf("error deleting shell user: %v", err)
	}

	if deleteRes.CallbackID != "" {
		if err = utils.WaitForAction(ctx, client, deleteRes.CallbackID); err != nil {
			return diag.Errorf("Error deleting shell user (%s): %v", d.Id(), err)
		}
	}

	return nil
}

func findShellUserByID(id string, shellUsers []webdock.ShellUser) *webdock.ShellUser {
	for _, shellUser := range shellUsers {
		if fmt.Sprintf("%d", shellUser.ID) == id {
			return &shellUser
		}
	}

	return nil
}

func setShellUserAttributes(d *schema.ResourceData, shellUser *webdock.ShellUser) error {
	d.SetId(fmt.Sprintf("%d", shellUser.ID))

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

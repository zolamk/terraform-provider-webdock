package webdock

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/api"
)

func resourceWebdockPublicKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebdockPublicKeyCreate,
		ReadContext:   resourceWebdockPublicKeyRead,
		DeleteContext: resourceWebdockPublicKeyDelete,
		SchemaVersion: 0,
		Schema:        publicKeySchema(),
	}
}

func resourceWebdockPublicKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	body := api.CreatePublicKeyModel{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("key").(string),
	}

	publicKey, err := client.CreatePublicKey(ctx, body)

	if err != nil {
		return diag.Errorf("error creating public key: %v", err)
	}

	if err = setPublicKeyAttributes(d, publicKey); err != nil {
		return diag.Errorf("error setting public key: %v", err)
	}

	return nil
}

func findPublicKeyById(id string, publicKeys *api.PublicKeys) *api.PublicKey {
	if publicKeys == nil {
		return nil
	}

	for _, publicKey := range *publicKeys {
		if publicKey.Id.String() == id {
			return &publicKey
		}
	}

	return nil
}

func setPublicKeyAttributes(d *schema.ResourceData, key *api.PublicKey) error {
	d.SetId(key.Id.String())

	if err := d.Set("name", key.Name); err != nil {
		return err
	}

	if err := d.Set("key", key.Key); err != nil {
		return err
	}

	if err := d.Set("created_at", key.Created); err != nil {
		return err
	}

	return nil
}

func resourceWebdockPublicKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	publicKeys, err := client.GetPublicKeys(ctx)

	if err != nil {
		return diag.Errorf("error getting public key: %v", err)
	}

	if err != nil {
		return diag.Errorf("error converting public key id to int64: %v", err)
	}

	publicKey := findPublicKeyById(d.Id(), publicKeys)

	if publicKey == nil {
		return diag.Errorf("error getting public key: 404 Not Found")
	}

	if err = setPublicKeyAttributes(d, publicKey); err != nil {
		return diag.Errorf("error setting public key: %v", err)
	}

	return nil
}

func resourceWebdockPublicKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	id, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.Errorf("error converting public key id to int64: %v", err)
	}

	if err = client.DeletePublicKey(ctx, id); err != nil {
		return diag.Errorf("error deleting public key: %v", err)
	}

	d.SetId("")

	return nil
}

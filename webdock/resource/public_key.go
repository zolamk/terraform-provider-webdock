package resource

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zolamk/terraform-provider-webdock/api"
	"github.com/zolamk/terraform-provider-webdock/config"
	"github.com/zolamk/terraform-provider-webdock/webdock/schemas"
)

func PublicKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: createPublicKey,
		ReadContext:   readPublicKey,
		DeleteContext: deletePublicKey,
		SchemaVersion: 0,
		Schema:        schemas.PublicKey(),
	}
}

func createPublicKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	body := api.CreatePublicKeyRequestBody{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("key").(string),
	}

	publicKey, err := client.CreatePublicKey(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = setPublicKeyAttributes(d, publicKey); err != nil {
		return diag.Errorf("error setting public key: %v", err)
	}

	return nil
}

func readPublicKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	publicKeys, err := client.GetPublicKeys(ctx)
	if err != nil {
		return diag.Errorf("error getting public key: %v", err)
	}

	publicKey := findPublicKeyById(d.Id(), publicKeys)

	if publicKey == nil {
		return diag.Errorf("error getting public key: not found")
	}

	if err = setPublicKeyAttributes(d, publicKey); err != nil {
		return diag.Errorf("error setting public key: %v", err)
	}

	return nil
}

func deletePublicKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	id, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.Errorf("error converting public key id to int64: %v", err)
	}

	if err = client.DeletePublicKey(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func findPublicKeyById(id string, publicKeys api.PublicKeys) *api.PublicKey {
	if publicKeys == nil {
		return nil
	}

	for _, publicKey := range publicKeys {
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

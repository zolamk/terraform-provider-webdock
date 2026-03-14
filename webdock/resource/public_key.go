package resource

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	webdock "github.com/webdock-io/go-sdk"
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

	opts := webdock.CreatePublicKeyOptions{
		Name:      d.Get("name").(string),
		PublicKey: d.Get("key").(string),
	}

	publicKey, err := client.CreatePublicKey(opts)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = setPublicKeyAttributes(d, &publicKey); err != nil {
		return diag.Errorf("error setting public key: %v", err)
	}

	return nil
}

func readPublicKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig)

	publicKeys, err := client.ListAccountPublicKeys(webdock.ListAccountPublicKeysOptions{})
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

	if err = client.DeletePublicKey(webdock.DeletePublicOptions{ID: id}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func findPublicKeyById(id string, publicKeys []webdock.PublicKey) *webdock.PublicKey {
	for _, publicKey := range publicKeys {
		if fmt.Sprintf("%d", publicKey.ID) == id {
			return &publicKey
		}
	}

	return nil
}

func setPublicKeyAttributes(d *schema.ResourceData, key *webdock.PublicKey) error {
	d.SetId(fmt.Sprintf("%d", key.ID))

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

package webdock

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceWebdockAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceWebdockAccountRead,
		Schema: map[string]*schema.Schema{
			"account_balance": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account credit balance display text",
			},
			"account_balance_raw": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account credit balance raw value",
			},
			"account_balance_raw_currency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account credit balance currency",
			},
			"company_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Company name",
			},
			"user_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User email",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User ID",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User name",
			},
		},
	}
}

func datasourceWebdockAccountRead(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).client

	account, err := client.GetAccountInformation(ctx)

	if err != nil {
		return diag.Errorf("Error getting account information: %s", err)
	}

	id := strconv.FormatInt(account.UserId, 10)

	rd.SetId(id)

	if err = rd.Set("account_balance", account.AccountBalance); err != nil {
		return diag.Errorf("Error setting account_balance: %s", err)
	}

	if err = rd.Set("account_balance_raw", account.AccountBalanceRaw); err != nil {
		return diag.Errorf("Error setting account_balance_raw: %s", err)
	}

	if err = rd.Set("account_balance_raw_currency", account.AccountBalanceRawCurrency); err != nil {
		return diag.Errorf("Error setting account_balance_raw_currency: %s", err)
	}

	if err = rd.Set("company_name", account.CompanyName); err != nil {
		return diag.Errorf("Error setting company_name: %s", err)
	}

	if err = rd.Set("user_email", account.UserEmail); err != nil {
		return diag.Errorf("Error setting user_email: %s", err)
	}

	if err = rd.Set("user_id", account.UserId); err != nil {
		return diag.Errorf("Error setting user_id: %s", err)
	}

	if err = rd.Set("user_name", account.UserName); err != nil {
		return diag.Errorf("Error setting user_name: %s", err)
	}

	return nil
}

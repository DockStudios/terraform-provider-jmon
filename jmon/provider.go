package jmon

import (
	//"log"
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"jmon_environment": resourceEnvironment(),
			"jmon_check":       resourceCheck(),
		},
	}
}

type ProviderClient struct {
	url        string
	apiKey     string
	httpClient *http.Client
	headers    http.Header
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := &ProviderClient{}

	url := d.Get("url").(string)
	if url == "" {
		url = "http://localhost:5000"
	}
	client.url = url

	client.apiKey = d.Get("api_key").(string)

	client.headers = make(http.Header)
	client.headers.Set("Content-Type", "application/json")
	client.headers.Set("Accept", "application/json")
	if client.apiKey != "" {
		client.headers.Set("X-JMon-Api-Key", client.apiKey)
	}

	client.httpClient = &http.Client{
		Timeout: time.Second * 30,
	}

	return client, diags
}

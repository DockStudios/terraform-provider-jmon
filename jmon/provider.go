package jmon

import (
	//"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ConfigureFunc: providerConfigure,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"jmon_check": resourceCheck(),
		},
	}
}

type ProviderClient struct {
	url        string
	httpClient *http.Client
	headers    http.Header
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &ProviderClient{}

	url := d.Get("url").(string)
	if url == "" {
		url = "http://localhost:5000"
	}
	client.url = url

	client.headers = make(http.Header)
	client.headers.Set("Content-Type", "application/json")
	client.headers.Set("Accept", "application/json")

	client.httpClient = &http.Client{
		Timeout: time.Second * 30,
	}

	return client, nil
}

package jmon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		// UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func generateEnvironmentUrl(client *ProviderClient, name string) string {
	return fmt.Sprintf("%s/api/v1/environments/%s", client.url, name)
}

func getEnvironmentByName(d *schema.ResourceData, m interface{}, name string, responseBody *[]byte) (error, bool) {
	client := m.(*ProviderClient)

	// Create request to GET environment
	var url string = generateEnvironmentUrl(client, name)
	log.Printf("[jmon] Perform GET: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, false
	}

	// Perform request
	r, err := client.httpClient.Do(req)
	if err != nil {
		return err, false
	}
	defer r.Body.Close()

	// Read response body
	*responseBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return err, false
	}

	log.Printf("[jmon] Response code: %d", r.StatusCode)
	log.Printf("[jmon] Body: %s", responseBody)

	var exists bool
	exists = true
	// Check status code to determine if check exists
	if r.StatusCode == 404 {
		exists = false
	}

	return nil, exists
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Determine if environment already exists
	var responseBody []byte
	err, exists := getEnvironmentByName(d, m, d.Get("name").(string), &responseBody)
	if err != nil {
		return diag.FromErr(err)
	}
	if exists {
		return diag.FromErr(errors.New(fmt.Sprintf("An environment already exists with the name: %s", d.Get("name").(string))))
	}

	// Create environment
	client := m.(*ProviderClient)
	var url string = fmt.Sprintf("%s/api/v1/environments", client.url)
	log.Printf("[jmon] Perform POST: %s", url)
	values := map[string]string{"name": d.Get("name").(string)}
	jsonPostData, err := json.Marshal(values)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPostData))
	if err != nil {
		return diag.FromErr(err)
	}

	// Perform request
	r, err := client.httpClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	// Read response body
	responseBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set ID of resource to environment name
	d.SetId(d.Get("name").(string))

	return diags
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	var responseBody []byte
	err, exists := getEnvironmentByName(d, m, d.Get("name").(string), &responseBody)

	if err != nil {
		return diag.FromErr(err)
	}

	// If check does not exist, reset ID
	// and return early
	if exists == false {
		log.Printf("[jmon] Environment does not exist: %s", string(responseBody))
		d.SetId("")
		return diags
	}

	// Set name attribute from ID of object.
	// This is important during imports, as the name attribute
	// does not yet exist, so will be updated from the imported ID
	d.Set("name", d.Id())

	return diags
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// There are no attributes that can be updated (yet),
	// so immediately return

	return diags
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*ProviderClient)

	// Create DELETE request
	req, err := http.NewRequest("DELETE", generateEnvironmentUrl(client, d.Get("name").(string)), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Perform request
	r, err := client.httpClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	// Read response body
	var responseBody []byte
	responseBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check status code
	if r.StatusCode != 200 {
		return diag.FromErr(errors.New(fmt.Sprintf("[jmon] Environment failed to delete: %s", string(responseBody))))
	}

	return diags
}

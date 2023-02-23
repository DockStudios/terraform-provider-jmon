package jmon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CheckData struct {
	Name              string        `yaml:"name"`
	Steps             []interface{} `yaml:"steps"`
	ScreenshotOnError bool          `yaml:"screenshot_on_error,omitempty"`
	Interval          int           `yaml:"interval,omitempty"`
	Client            string        `yaml:"client,omitempty"`
}

func resourceCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCheckCreate,
		ReadContext:   resourceCheckRead,
		UpdateContext: resourceCheckUpdate,
		DeleteContext: resourceCheckDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"steps": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"screenshot_on_error": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"client": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func upsertCheck(d *schema.ResourceData, m interface{}, check *CheckData) error {

	// Convert resource attributes to attributes of check data object
	check.Name = d.Get("name").(string)
	check.Interval = d.Get("interval").(int)
	check.ScreenshotOnError = d.Get("screenshot_on_error").(bool)
	check.Client = d.Get("client").(string)

	// Convert steps YAML to interface in check object
	ymlErr := yaml.Unmarshal([]byte(d.Get("steps").(string)), &check.Steps)
	if ymlErr != nil {
		return ymlErr
	}

	client := m.(*ProviderClient)

	// Convert entire check to YAML
	var yamlOutput []byte
	yamlOutput, err := yaml.Marshal(&check)
	if err != nil {
		return err
	}

	// Create reader for post data
	postDataReader := bytes.NewReader(yamlOutput)

	// Check request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/checks", client.url), postDataReader)
	if err != nil {
		return err
	}

	// Perform request
	r, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Read response body
	var responseBody []byte
	responseBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Check status code
	if r.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Check failed to create/update: %s", string(responseBody)))
	}

	return nil
}

func resourceCheckCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Determine if check already exists
	var responseBody []byte
	err, exists := getCheckByName(d, m, d.Get("name").(string), &responseBody)
	if err != nil {
		return diag.FromErr(err)
	}
	if exists {
		return diag.FromErr(errors.New(fmt.Sprintf("A check already exists with the name: %s", d.Get("name").(string))))
	}

	var check CheckData
	err = upsertCheck(d, m, &check)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set ID of resource to check name
	d.SetId(check.Name)

	return diags
}

func getCheckByName(d *schema.ResourceData, m interface{}, name string, responseBody *[]byte) (error, bool) {
	client := m.(*ProviderClient)

	// Check request
	log.Printf("[jmon] Perform GET: %s/api/v1/checks/%s", client.url, name)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/checks/%s", client.url, name), nil)
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

func resourceCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	var responseBody []byte
	err, exists := getCheckByName(d, m, d.Id(), &responseBody)

	if err != nil {
		return diag.FromErr(err)
	}

	// If check does not exist, reset ID
	// and return early
	if exists == false {
		log.Printf("[jmon] Check does not exist: %s", string(responseBody))
		d.SetId("")
		return diags
	}

	var check CheckData
	err = yaml.Unmarshal(responseBody, &check)
	log.Printf("[jmon] Unmarshalled to: %v", check)
	if err != nil {
		return diag.FromErr(err)
	}

	stepsString, err := yaml.Marshal(&check.Steps)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set name attribute from ID of object.
	// This is important during imports, as the name attribute
	// does not yet exist, so will be updated from the imported ID
	d.Set("name", d.Id())

	d.Set("interval", check.Interval)
	d.Set("client", check.Client)
	d.Set("steps", string(stepsString))
	d.Set("screenshot_on_error", check.ScreenshotOnError)

	return diags
}

func resourceCheckUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	var check CheckData
	err := upsertCheck(d, m, &check)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCheckDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*ProviderClient)

	// Check request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/checks/%s", client.url, d.Get("name").(string)), nil)
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
		return diag.FromErr(errors.New(fmt.Sprintf("[jmon] Check failed to delete: %s", string(responseBody))))
	}

	return diags
}

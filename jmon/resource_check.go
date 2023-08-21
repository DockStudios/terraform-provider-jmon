package jmon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CheckData struct {
	Name              string                 `yaml:"name"`
	Environment       string                 `yaml:"environment,omitempty"`
	Steps             []interface{}          `yaml:"steps"`
	ScreenshotOnError bool                   `yaml:"screenshot_on_error,omitempty"`
	Interval          int                    `yaml:"interval,omitempty"`
	Timeout           int                    `yaml:"timeout,omitempty"`
	Client            string                 `yaml:"client,omitempty"`
	Enable            bool                   `yaml:"enable,omitempty"`
	Attributes        map[string]interface{} `yaml:"attributes,omitempty"`
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
			"environment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"steps": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"screenshot_on_error": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"client": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func generateId(name string, environment string) string {
	var id string
	id = name
	if environment != "" {
		id = id + "/" + environment
	}
	return id
}

func nameEnvironmentFromId(id string) (string, string, error) {
	splitNames := strings.Split(id, "/")
	if len(splitNames) == 1 {
		return splitNames[0], "", nil
	} else if len(splitNames) == 2 {
		return splitNames[0], splitNames[1], nil
	} else {
		return "", "", errors.New("Cannot parse invalid ID")
	}
}

func generateCheckUrl(client *ProviderClient, name string, environment string) string {
	return fmt.Sprintf("%s/api/v1/checks/%s/environments/%s", client.url, name, environment)
}

func upsertCheck(d *schema.ResourceData, m interface{}, check *CheckData) error {

	// Convert resource attributes to attributes of check data object
	check.Name = d.Get("name").(string)
	check.Environment = d.Get("environment").(string)
	check.Interval = d.Get("interval").(int)
	check.Timeout = d.Get("timeout").(int)
	check.ScreenshotOnError = d.Get("screenshot_on_error").(bool)
	check.Client = d.Get("client").(string)
	check.Enable = d.Get("enable").(bool)
	check.Attributes = d.Get("attributes").(map[string]interface{})

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

	// Add headers to request
	req.Header = client.headers.Clone()

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
	err, exists := getCheckByNameAndEnvironment(d, m, d.Get("name").(string), d.Get("environment").(string), &responseBody)
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
	d.SetId(generateId(check.Name, check.Environment))

	return diags
}

func getCheckByNameAndEnvironment(d *schema.ResourceData, m interface{}, name string, environment string, responseBody *[]byte) (error, bool) {
	client := m.(*ProviderClient)

	// Check request
	var url string = generateCheckUrl(client, name, environment)
	log.Printf("[jmon] Perform GET: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, false
	}

	// Add headers to request
	req.Header = client.headers.Clone()

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

	name, environment, err := nameEnvironmentFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	// Set default environment, if it does not exist
	if environment == "" {
		environment = "default"
		d.SetId(generateId(name, "default"))
	}

	var responseBody []byte
	err, exists := getCheckByNameAndEnvironment(d, m, name, environment, &responseBody)

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

	// Only set steps string from API response, if the strctures differ.
	// If they aren't different, do not update, comments and insignificant
	// whitespace changes will cause a different in value, even if the
	// generated structure is fundermentally the same.
	updateSteps := true

	// Unmarshall current steps and perform deep compare
	var currentSteps []interface{}
	ymlErr := yaml.Unmarshal([]byte(d.Get("steps").(string)), &currentSteps)
	if ymlErr == nil {
		if reflect.DeepEqual(currentSteps, check.Steps) {
			log.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!! STEPS ARE EQUALS")
			updateSteps = false
		}
	}

	if updateSteps {
		stepsString, err := yaml.Marshal(&check.Steps)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("steps", string(stepsString))
	}

	// Set name attribute from ID of object.
	// This is important during imports, as the name attribute
	// does not yet exist, so will be updated from the imported ID
	d.Set("name", name)
	d.Set("environment", environment)

	d.Set("interval", check.Interval)
	d.Set("timeout", check.Timeout)
	d.Set("client", check.Client)
	d.Set("screenshot_on_error", check.ScreenshotOnError)
	d.Set("enable", check.Enable)
	d.Set("attributes", check.Attributes)

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
	req, err := http.NewRequest("DELETE", generateCheckUrl(client, d.Get("name").(string), d.Get("environment").(string)), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Add headers to request
	req.Header = client.headers.Clone()

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

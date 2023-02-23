package jmon

import (
	"bytes"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type CheckData struct {
	name                string        `yaml:"name"`
	steps               []interface{} `yaml:"steps"`
	screenshot_on_error bool          `yaml:"screenshot_on_error"`
	interval            int           `yaml:"interval"`
	client              string        `yaml:"client"`
}

func resourceCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceCheckCreate,
		Read:   resourceCheckRead,
		Update: resourceCheckUpdate,
		Delete: resourceCheckDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func resourceCheckCreate(d *schema.ResourceData, m interface{}) error {

	var check CheckData

	// Convert resource attributes to attributes of check data object
	check.name = d.Get("name").(string)
	check.interval = d.Get("interval").(int)
	check.screenshot_on_error = d.Get("screenshot_on_error").(bool)
	check.client = d.Get("client").(string)

	// Convert steps YAML to interface in check object
	ymlErr := yaml.Unmarshal([]byte(d.Get("steps").(string)), &check.steps)
	if ymlErr != nil {
		return ymlErr
	}

	client := m.(*ProviderClient)

	// Convert entire check to YAML
	var yamlOutput []byte
	yamlOutput, err := yaml.Marshal(check)

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

	// Set ID of resource to check name
	d.SetId(check.name)

	return nil
}

func resourceCheckRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCheckUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCheckDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

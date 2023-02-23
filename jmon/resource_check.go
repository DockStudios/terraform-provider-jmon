package jmon

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	fmt.Printf("YAML: %s", yamlOutput)

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
	_, err = r.Body.Read(responseBody)
	if err != nil {
		return err
	}

	// Check status code
	if r.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Check failed to create: %s", string(responseBody[:])))
	}

	return nil
}

func resourceCheckCreate(d *schema.ResourceData, m interface{}) error {

	var check CheckData
	err := upsertCheck(d, m, &check)

	if err != nil {
		return err
	}

	// Set ID of resource to check name
	d.SetId(check.Name)

	return nil
}

func resourceCheckRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCheckUpdate(d *schema.ResourceData, m interface{}) error {

	var check CheckData
	err := upsertCheck(d, m, &check)

	if err != nil {
		return err
	}

	return nil
}

func resourceCheckDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

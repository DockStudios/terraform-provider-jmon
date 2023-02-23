package jmon

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

type JmonResponse struct {
	message string `yaml:"msg"`
	status  string `yaml:"status"`
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
	responseBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Check status code
	if r.StatusCode != 200 {
		var response JmonResponse
		getResponseFromBody(&responseBody, &response)
		return errors.New(fmt.Sprintf("Check failed to create/update: %s", response.message))
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

func getResponseFromBody(responseBody *[]byte, response *JmonResponse) bool {
	err := yaml.Unmarshal(*responseBody, response)
	log.Printf("[jmon] Unmarshalled response to: %v", response)
	if err != nil {
		return false
	}
	return true
}

func getCheckByName(d *schema.ResourceData, m interface{}, responseBody *[]byte) (error, bool) {
	client := m.(*ProviderClient)

	// Check request
	log.Printf("[jmon] Perform GET: %s/api/v1/checks/%s", client.url, d.Get("name").(string))
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/checks/%s", client.url, d.Get("name").(string)), nil)
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
	//_, err = r.Body.Read(responseBody)
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

func resourceCheckRead(d *schema.ResourceData, m interface{}) error {

	var responseBody []byte
	err, exists := getCheckByName(d, m, &responseBody)

	if err != nil {
		return err
	}

	// If check does not exist, reset ID
	// and return early
	if exists == false {
		log.Printf("[jmon] Check does not exist: %s", string(responseBody))
		d.SetId("")
		return nil
	}

	var check CheckData
	err = yaml.Unmarshal(responseBody, &check)
	log.Printf("[jmon] Unmarshalled to: %v", check)
	if err != nil {
		return err
	}

	stepsString, err := yaml.Marshal(&check.Steps)
	if err != nil {
		return err
	}

	d.Set("interval", check.Interval)
	d.Set("client", check.Client)
	d.Set("steps", string(stepsString))
	d.Set("screenshot_on_error", check.ScreenshotOnError)

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
	client := m.(*ProviderClient)

	// Check request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/checks/%s", client.url, d.Get("name").(string)), nil)
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
		var response JmonResponse
		getResponseFromBody(&responseBody, &response)
		return errors.New(fmt.Sprintf("Check failed to delete: %s", response.message))
	}

	return nil
}

package authfed

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"net/http"
	"strings"
)

func resourceHttpObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceHttpObjectCreate,
		Read:   resourceHttpObjectRead,
		Update: resourceHttpObjectUpdate,
		Delete: resourceHttpObjectDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceHttpObjectCreate(d *schema.ResourceData, m interface{}) error {
	return resourceHttpObjectUpdate(d, m)
}

func resourceHttpObjectRead(d *schema.ResourceData, m interface{}) error {
	client := getClient(m)
	url := d.Get("url").(string)
	if url == "" {
		url = d.Id()
		d.Set("url", d.Id())
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("request error: %s (%s)", url, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not GET object: %s (%s)", url, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read body of object: %s (%s)", url, err)
	}
	if err = d.Set("content", string(body)); err != nil {
		return fmt.Errorf("error while setting 'content': %s", err)
	}
	d.SetId(url)
	return nil
}

func resourceHttpObjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := getClient(m)
	url := d.Get("url").(string)
	content := d.Get("content").(string)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("request error: %s (%s)", url, err)
	}
	if _, err = client.Do(req); err != nil {
		return fmt.Errorf("could not PUT object: %s (%s)", url, err)
	}
	return resourceHttpObjectRead(d, m)
}

func resourceHttpObjectDelete(d *schema.ResourceData, m interface{}) error {
	client := getClient(m)
	url := d.Get("url").(string)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("request error: %s (%s)", url, err)
	}
	if _, err = client.Do(req); err != nil {
		return fmt.Errorf("could not DELETE object: %s (%s)", url, err)
	}
	return nil
}

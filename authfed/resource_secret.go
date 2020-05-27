package authfed

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"os"
	"strings"
)

func resourceHttpSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceHttpSecretCreate,
		Read:   resourceHttpSecretRead,
		Update: resourceHttpSecretUpdate,
		Delete: resourceHttpSecretDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"filename": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"checksum": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceHttpSecretCreate(d *schema.ResourceData, m interface{}) error {
	return resourceHttpSecretUpdate(d, m)
}

func resourceHttpSecretRead(d *schema.ResourceData, m interface{}) error {
	client := getClient(m)
	url := d.Get("url").(string)
	if url == "" {
		url = d.Id()
		d.Set("url", d.Id())
	}
	checksum := d.Get("checksum").(string)
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return fmt.Errorf("request error: %s (%s)", url, err)
	}
	req.Header.Set("If-None-Match", checksum)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not HEAD object: %s (%s)", url, err)
	}
	etag := resp.Header.Get("ETag")
	if etag == "" {
		return fmt.Errorf("could not get ETag of object: %s (%s)", url, etag)
	}
	if err = d.Set("checksum", strings.ReplaceAll(etag, "\"", "")); err != nil {
		return fmt.Errorf("error while setting 'checksum': %s", err)
	}
	d.SetId(url)
	d.Set("filename", nil)
	return nil
}

func resourceHttpSecretUpdate(d *schema.ResourceData, m interface{}) error {
	client := getClient(m)
	url := d.Get("url").(string)
	filename := d.Get("filename").(string)
	data, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("file open error: %s (%s)", url, err)
	}
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		return fmt.Errorf("request error: %s (%s)", url, err)
	}
	if _, err = client.Do(req); err != nil {
		return fmt.Errorf("could not PUT object: %s (%s)", url, err)
	}
	return resourceHttpSecretRead(d, m)
}

func resourceHttpSecretDelete(d *schema.ResourceData, m interface{}) error {
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

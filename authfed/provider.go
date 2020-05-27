package authfed

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/mapstructure"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cert_file": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_file": {
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"authfed_http_object": resourceHttpObject(),
			"authfed_http_secret": resourceHttpSecret(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var config *Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func getClient(meta interface{}) *http.Client {
	// We can ignore err since we checked the configuration in providerConfigure()
	client, _ := meta.(*Config).Client()
	return client
}

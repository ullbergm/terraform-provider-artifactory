package artifactory

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/webdevwilson/go-artifactory/artifactory"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// Provider returns a terraform.resourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARTIFACTORY_USERNAME", nil),
				Description: "Username for authentication",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARTIFACTORY_PASSWORD", nil),
				Description: "Password or API Key to use",
			},

			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARTIFACTORY_URL", nil),
				Description: "The URL to your Artifactory instance ",
			},

			"skip_tls_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARTIFACTORY_SKIP_VERIFY", ""),
				Description: "Set this to true only if the target Artifactory server is an insecure development instance.",
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"artifactory_local_repository":   resourceLocalRepository(),
			"artifactory_remote_repository":  resourceRemoteRepository(),
			"artifactory_virtual_repository": resourceVirtualRepository(),
			"artifactory_user":               resourceUser(),
			"artifactory_group":              resourceGroup(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	user := d.Get("username").(string)
	pass := d.Get("password").(string)
	url := d.Get("url").(string)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: d.Get("skip_tls_verify").(bool)}
	hc := &http.Client{Transport: http.DefaultTransport}
	hc.Transport = logging.NewTransport("Artifactory", hc.Transport)
	c := artifactory.NewClient(user, pass, url, hc)

	// fail early. validate the connection to Artifactory
	if err := c.Ping(); err != nil {
		return nil, fmt.Errorf("Error connecting to Artifactory: %s", err)
	}

	return c, nil
}

package ipsfor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

var namikodaUrl = "https://api.namikoda.com"

func dataSource() *schema.Resource {
	return innerDataSource("")
}

func innerDataSource(url string) *schema.Resource {
	if len(url) > 0 {
		namikodaUrl = url
	}
	return &schema.Resource{
		Read: dataSourceRead,

		Schema: map[string]*schema.Schema{
			"apikey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ipv4s": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ipv6s": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"value": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"lastUpdate": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(d *schema.ResourceData, meta interface{}) error {

	apikey := d.Get("apikey").(string)
	id := d.Get("id").(string)

	owner := "public"
	ownerFromConfig, ownerOk := d.GetOk("owner")
	if ownerOk {
		owner = ownerFromConfig.(string)
	}

	client := &http.Client{}

	path := url.PathEscape(id)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/%s/ipsfor/%s", namikodaUrl, owner, path), nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}

	req.Header.Set("X-Namikoda-Key", apikey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error %s during making a request: %s", err, "https://api.namikoda.com/v1/public/ipsfor/"+path)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || isContentTypeAllowed(contentType) == false {
		return fmt.Errorf("Content-Type is not a text type. Got: %s", contentType)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading response body. %s", err)
	}

	ipsfor := new(Ipsfor)
	err = json.Unmarshal(bytes, &ipsfor)
	if err != nil {
		return fmt.Errorf("Error while parsing response body. %s", bytes)
	}
	d.Set("ipv4s", ipsfor.Ipv4s)
	d.Set("ipv6s", ipsfor.Ipv6s)
	d.Set("value", ipsfor.Value)
	d.Set("lastUpdate", ipsfor.LastUpdate)
	d.SetId("namikoda-" + id + "-" + ipsfor.LastUpdate)

	return nil
}

// IpsFor
//https://mholt.github.io/json-to-go/
type Ipsfor struct {
	Ipv4s      []string `json:"ipv4s"`
	Ipv6s      []string `json:"ipv6s"`
	LastUpdate string   `json:"lastUpdate"`
	Name       string   `json:"name"`
	ID         string   `json:"id"`
	Value      []string `json:"value"`
}

// This is to prevent potential issues w/ binary files
// and generally unprintable characters
// See https://github.com/hashicorp/terraform/pull/3858#issuecomment-156856738
func isContentTypeAllowed(contentType string) bool {
	allowedContentTypes := []*regexp.Regexp{
		regexp.MustCompile("^application/json$"),
	}

	for _, r := range allowedContentTypes {
		if r.MatchString(contentType) {
			return true
		}
	}

	return false
}

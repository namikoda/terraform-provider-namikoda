package ipsfor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

type TestHttpMock struct {
	server *httptest.Server
}

func createTestProvider(url string) map[string]terraform.ResourceProvider {
	return map[string]terraform.ResourceProvider{
		"ipsfor": &schema.Provider{
			Schema: map[string]*schema.Schema{},

			DataSourcesMap: map[string]*schema.Resource{
				"ipsfor": innerDataSource(url),
			},

			ResourcesMap: map[string]*schema.Resource{},
		},
	}
}

// reflect.DeepEqual has problems with []interface{} vs []string
func compareExpectation(a interface{}, b []string) bool {
	//if len(a) != len(b) {
	//	return false
	//}
	//for i := 0; i < len(a); i++ {
	//	if a[i] != b[i] {
	//		return false
	//	}
	//}
	return true
}

const testDataSourceConfig_basic = `
data "ipsfor" "dummy-success-test" {
  apikey = "aaaabbbb-cccc-dddd-eeee-ffffgggghhhh"
  id = "dummy-success"
}

output "ipv4s" {
  value = "${data.ipsfor.dummy-success-test.ipv4s}"
}

output "ipv6s" {
  value = "${data.ipsfor.dummy-success-test.ipv6s}"
}

output "value" {
  value = "${data.ipsfor.dummy-success-test.value}"
}

output "lastUpdate" {
  value = "${data.ipsfor.dummy-success-test.lastUpdate}"
}
`

func TestDataSource_basic(t *testing.T) {
	testHttpMock := setUpMockHttpServer()
	testProviders = createTestProvider(testHttpMock.server.URL)
	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testDataSourceConfig_basic,
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.ipsfor.dummy-success-test"]
					if !ok {
						printable, _ := json.Marshal(s.RootModule().Resources)
						return fmt.Errorf("missing ipsfor data resource: %s", printable)
					}

					outputs := s.RootModule().Outputs

					ipv4s := make([]interface{}, 1)
					ipv4s[0] = "1.2.3.4/32"
					if !reflect.DeepEqual(outputs["ipv4s"].Value, ipv4s) {
						return fmt.Errorf(`'ipv4s' output is %s; want %s`, outputs["ipv4s"].Value, ipv4s)
					}

					ipv6s := make([]interface{}, 1)
					ipv6s[0] = "1111:222:3000::/44"
					if !reflect.DeepEqual(outputs["ipv6s"].Value, ipv6s) {
						return fmt.Errorf(`'ipv6s' output is %s; want %s`, outputs["ipv6s"].Value, ipv6s)
					}

					value := make([]interface{}, 2)
					value[0] = "1.2.3.4/32"
					value[1] = "1111:222:3000::/44"
					if !reflect.DeepEqual(outputs["value"].Value, value) {
						return fmt.Errorf(`'value' output is %s; want %s`, outputs["value"].Value, value)
					}

					if outputs["lastUpdate"].Value != "2017-01-01T00:00:00.000Z" {
						return fmt.Errorf(
							`'lastUpdate' output is %s; want "2017-01-01T00:00:00.000Z"`,
							outputs["lastUpdate"].Value,
						)
					}

					return nil
				},
			},
		},
	})
}

const testDataSourceConfig_otherowner = `
data "ipsfor" "dummy-success-test" {
  apikey = "aaaabbbb-cccc-dddd-eeee-ffffgggghhhh"
  owner = "otherowner"
  id = "dummy-success"
}

output "ipv4s" {
  value = "${data.ipsfor.dummy-success-test.ipv4s}"
}

output "ipv6s" {
  value = "${data.ipsfor.dummy-success-test.ipv6s}"
}

output "value" {
  value = "${data.ipsfor.dummy-success-test.value}"
}

output "lastUpdate" {
  value = "${data.ipsfor.dummy-success-test.lastUpdate}"
}
`

func TestDataSource_otherowner(t *testing.T) {
	testHttpMock := setUpMockHttpServer()
	testProviders = createTestProvider(testHttpMock.server.URL)
	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testDataSourceConfig_otherowner,
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.ipsfor.dummy-success-test"]
					if !ok {
						printable, _ := json.Marshal(s.RootModule().Resources)
						return fmt.Errorf("missing ipsfor data resource: %s", printable)
					}

					outputs := s.RootModule().Outputs

					ipv4s := make([]interface{}, 1)
					ipv4s[0] = "5.6.7.8/32"
					if !reflect.DeepEqual(outputs["ipv4s"].Value, ipv4s) {
						return fmt.Errorf(`'ipv4s' output is %s; want %s`, outputs["ipv4s"].Value, ipv4s)
					}

					ipv6s := make([]interface{}, 1)
					ipv6s[0] = "5555:666:7000::/44"
					if !reflect.DeepEqual(outputs["ipv6s"].Value, ipv6s) {
						return fmt.Errorf(`'ipv6s' output is %s; want %s`, outputs["ipv6s"].Value, ipv6s)
					}

					value := make([]interface{}, 2)
					value[0] = "5555:666:7000::/44"
					value[1] = "5.6.7.8/32"
					if !reflect.DeepEqual(outputs["value"].Value, value) {
						return fmt.Errorf(`'value' output is %s; want %s`, outputs["value"].Value, value)
					}

					if outputs["lastUpdate"].Value != "2017-01-01T00:00:00.000Z" {
						return fmt.Errorf(
							`'lastUpdate' output is %s; want "2017-01-01T00:00:00.000Z"`,
							outputs["lastUpdate"].Value,
						)
					}

					return nil
				},
			},
		},
	})
}

const testDataSourceUnauthorized_error = `
data "ipsfor" "dummy-failure-unauthorized" {
  apikey = "nope"
  id = "dummy-failure"
}
`

func TestDataSourceUnauthorized_error(t *testing.T) {
	testHttpMock := setUpMockHttpServer()
	testProviders = createTestProvider(testHttpMock.server.URL)
	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testDataSourceUnauthorized_error,
				ExpectError: regexp.MustCompile("HTTP request error. Response code: 401"),
			},
		},
	})
}

const testDataSourceUnfound_error = `
data "ipsfor" "dummy-failure-unfound" {
  apikey = "aaaabbbb-cccc-dddd-eeee-ffffgggghhhh"
  id = "dummy-failure"
}
`

func TestDataSourceUnfound_error(t *testing.T) {
	testHttpMock := setUpMockHttpServer()
	testProviders = createTestProvider(testHttpMock.server.URL)
	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testDataSourceUnfound_error,
				ExpectError: regexp.MustCompile("HTTP request error. Response code: 404"),
			},
		},
	})
}

const testDataSourceConfig_error = `
data "ipsfor" "ipsfor-tooempty" {

}
`

func TestDataSource_compileError(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testDataSourceConfig_error,
				ExpectError: regexp.MustCompile("required field is not set"),
			},
		},
	})
}

func setUpMockHttpServer() *TestHttpMock {
	Server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Namikoda-Key") == "aaaabbbb-cccc-dddd-eeee-ffffgggghhhh" {
				if r.URL.Path == "/v1/public/ipsfor/dummy-success" {
					w.Header().Add("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{ "ipv4s":[ "1.2.3.4/32" ], "ipv6s":[ "1111:222:3000::/44" ], "lastUpdate":"2017-01-01T00:00:00.000Z", "name":"Dummy success value", "owner":"public", "id":"dummy-success", "value":[ "1.2.3.4/32", "1111:222:3000::/44" ] }`))
				} else if r.URL.Path == "/v1/otherowner/ipsfor/dummy-success" {
					w.Header().Add("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{ "ipv4s":[ "5.6.7.8/32" ], "ipv6s":[ "5555:666:7000::/44" ], "lastUpdate":"2017-01-01T00:00:00.000Z", "name":"Dummy success value", "owner":"otherowner", "id":"dummy-success", "value":[ "5555:666:7000::/44", "5.6.7.8/32" ] }`))
				} else {
					http.Error(w, "Not found: "+r.Header.Get("X-Namikoda-Key"), http.StatusNotFound)
				}
			} else {
				http.Error(w, "Invalid auth header: "+r.Header.Get("X-Namikoda-Key"), http.StatusUnauthorized)
			}
		}),
	)

	return &TestHttpMock{
		server: Server,
	}
}

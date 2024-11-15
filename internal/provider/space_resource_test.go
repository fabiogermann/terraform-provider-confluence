package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"terraform-provider-confluence/internal/fakeserver"
	"terraform-provider-confluence/internal/provider/transferobjects"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func generateTestSpaceObject() transferobjects.Space {
	base := transferobjects.Space{
		Id:   123,
		Name: "asdasd",
		Key:  "KEY",
	}
	return base
}

func TestSpaceResourceResource(t *testing.T) {
	t.SkipNow()
	debug := true
	apiServerObjects := make(map[string]map[string]interface{})

	svr := fakeserver.NewFakeServer(testPost, apiServerObjects, true, debug, "")
	test_url := fmt.Sprintf(`http://%s:%d`, testHost, testPost)
	os.Setenv("REST_API_URI", test_url)

	opt := &fakeserver.ApiClientOpt{
		Uri:                 test_url,
		Insecure:            false,
		Username:            "",
		Password:            "",
		Headers:             make(map[string]string),
		Timeout:             200,
		IdAttribute:         "id",
		CopyKeys:            make([]string, 0),
		WriteReturnsObject:  false,
		CreateReturnsObject: false,
		Debug:               debug,
	}
	client, err := fakeserver.NewAPIClient(opt)
	if err != nil {
		t.Fatal(err)
	}

	svr.SetSplice("/rest/api/space", func(a string, b []byte) (string, map[string]interface{}) {
		id := generateTestSpaceObject().Id
		jsonStr, _ := json.Marshal(generateTestSpaceObject())
		var obj map[string]interface{}
		_ = json.Unmarshal(jsonStr, &obj)
		return id.String(), obj
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			svr.StartInBackground()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExceptionItemResourceConfig(generateTestSpaceObject(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					fakeserver.TestAccCheckRestapiObjectExists("confluence_space.test", "id", client),
					resource.TestCheckResourceAttr("confluence_space.test", "key", generateTestSpaceObject().Key),
				),
			},
			// ImportState testing
			{
				ResourceName:      "confluence_space.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"key", "name"},
			},
			// Update and Read testing
			{
				Config: testAccExceptionItemResourceConfig(generateTestSpaceObject(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("confluence_space.test", "key", generateTestSpaceObject().Key),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

	svr.Shutdown()
}

func testAccExceptionItemResourceConfig(space transferobjects.Space, name string) string {
	return fmt.Sprintf(`%s
resource "confluence_space" "%s" {
  key = "%s"
  name = "%s"
}
`, providerConfig, name, space.Key, space.Name)
}

package provider

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"terraform-provider-confluence/internal/fakeserver"
	"terraform-provider-confluence/internal/provider/transferobjects"
	"testing"
)

func generateTestResource() transferobjects.Group {
	ruleContent := transferobjects.Group{
		Name: "groupName",
		Id:   "testId",
	}
	return ruleContent
}

func TestAccGroupResource(t *testing.T) {
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
		Timeout:             2,
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

	svr.SetSplice("/rest/api/group", func(a string, b []byte) (string, map[string]interface{}) {
		id := generateTestResource().Id
		var obj map[string]interface{}
		obj = structs.Map(generateTestResource())
		return id, obj
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
				Config: testAccDetectionRuleResourceConfig(generateTestResource(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					fakeserver.TestAccCheckRestapiObjectExists("confluence_group.test", "id", client),
					resource.TestCheckResourceAttr("confluence_group.test", "name", generateTestResource().Name),
				),
			},
			// ImportState testing
			{
				ResourceName:      "confluence_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"name"},
			},
			// Update and Read testing
			{
				Config: testAccDetectionRuleResourceConfig(generateTestResource(), "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("confluence_group.test", "name", generateTestResource().Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

	svr.Shutdown()
}

func testAccDetectionRuleResourceConfig(group transferobjects.Group, name string) string {
	return fmt.Sprintf(`%s
resource "confluence_group" "%s" {
 name = "%s"
}
`, providerConfig, name, group.Name)
}

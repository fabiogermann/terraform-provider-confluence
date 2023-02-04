package provider

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"strings"
	"terraform-provider-confluence/internal/fakeserver"
	"terraform-provider-confluence/internal/provider/transferobjects"
	"testing"
)

func generateTestSpacePermission() (string, string, []string) {
	key := "KEY"
	group := "groupName"
	permissions := []string{"create:page", "create:blogpost", "create:comment", "create:attachment"}
	return key, group, permissions
}

func TestAccExceptionContainerResource(t *testing.T) {
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

	key, group, permission := generateTestSpacePermission()

	path := fmt.Sprintf("/rest/api/space/%s", key)

	svr.SetSplice(path, func(a string, b []byte) (string, map[string]interface{}) {
		summary := testAccGenerateSpacePermissionObjects(key, group, permission)
		var obj map[string]interface{}
		obj = structs.Map(summary)
		return summary.ID.String(), obj
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
				Config: testAccSpacePermissionResourceConfig(key, group, permission, "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					fakeserver.TestAccCheckRestapiObjectExists("confluence_space_permission.test", "id", client),
					resource.TestCheckResourceAttr("confluence_space_permission.test", "key", key),
				),
			},
			// ImportState testing
			//{
			//	ResourceName:      "confluence_space_permission.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//	// This is not normally necessary, but is here because this
			//	// example code does not have an actual upstream service.
			//	// Once the Read method is able to refresh information from
			//	// the upstream service, this can be removed.
			//	ImportStateVerifyIgnore: []string{"rule_content"},
			//},
			// Update and Read testing
			{
				Config: testAccSpacePermissionResourceConfig(key, group, permission, "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("confluence_space_permission.test", "key", key),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})

	svr.Shutdown()
}

func testAccSpacePermissionResourceConfig(key string, group string, permissions []string, name string) string {
	return fmt.Sprintf(`%s
resource "confluence_space_permission" "%s" {
 key = "%s"
 group = "%s"
 operations = ["%s"]
}
`, providerConfig, name,
		key,
		group,
		strings.Join(permissions, "\", \""),
	)
}

func testAccGenerateSpacePermissionObjects(key string, group string, permissions []string) transferobjects.SummarySpacePermissions {
	var permissionObjects []transferobjects.SavedPermission
	for _, item := range permissions {
		var gresults []transferobjects.SavedPermissionGroupResult
		gresults = append(gresults, transferobjects.SavedPermissionGroupResult{
			Type: "group",
			Name: group,
			ID:   group,
		})
		groupp := transferobjects.SavedPermissionGroup{
			Results: gresults,
		}
		subjects := transferobjects.SavedPermissionSubjects{
			Group: &groupp,
		}
		permissionParts := strings.Split(item, ":")
		operation := transferobjects.SavedPermissionOperation{
			Operation:  permissionParts[0],
			TargetType: permissionParts[1],
		}
		var permissionObject = transferobjects.SavedPermission{
			ID:        1234,
			Subjects:  subjects,
			Operation: operation,
		}
		permissionObjects = append(permissionObjects, permissionObject)
	}

	summary := transferobjects.SummarySpacePermissions{
		ID:          123,
		Key:         key,
		Name:        "name",
		Type:        "type",
		Status:      "status",
		Permissions: permissionObjects,
	}
	return summary
}

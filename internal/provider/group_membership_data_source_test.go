package provider

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"terraform-provider-confluence/internal/fakeserver"
	"terraform-provider-confluence/internal/helpers"
	"testing"
)

var (
	testGroupId = "testid"
)

func TestAccPrivilegesDataSource(t *testing.T) {
	debug := true
	apiServerObjects := make(map[string]map[string]interface{})
	testSpaceObject := generateTestSpaceObject()

	svr := fakeserver.NewFakeServer(testPost, apiServerObjects, true, debug, "")
	test_url := fmt.Sprintf(`http://%s:%d`, testHost, testPost)
	os.Setenv("REST_API_URI", test_url)

	path := fmt.Sprintf("/rest/api/group/%s/membersByGroupId", testGroupId)
	setSliceA(svr, path, testSpaceObject.Id.String())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			svr.StartInBackground()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPrivilegesDataSourceConfig("test", testGroupId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.confluence_group_membership.test", "id", helpers.Sha256String(testGroupId)),
				),
			},
		},
	})
}

func testAccPrivilegesDataSourceConfig(name string, groupId string) string {
	return fmt.Sprintf(`%s
data "confluence_group_membership" "%s" {
	group_id = "%s"
}
`, providerConfig, name, groupId)
}

func setSliceA(svr *fakeserver.Fakeserver, path string, id string) {
	svr.SetSplice(path, func(a string, b []byte) (string, map[string]interface{}) {
		var obj map[string]interface{}
		_ = json.Unmarshal([]byte(generateGroupMembersResponseJson(0, 200, 200, 201)), &obj)
		setSliceB(svr, path, id)
		return id, obj
	})
}
func setSliceB(svr *fakeserver.Fakeserver, path string, id string) {
	svr.SetSplice(path, func(a string, b []byte) (string, map[string]interface{}) {
		var obj map[string]interface{}
		_ = json.Unmarshal([]byte(generateGroupMembersResponseJson(200, 200, 1, 201)), &obj)
		setSliceA(svr, path, id)
		return id, obj
	})
}

func generateGroupMembersResponseJson(start int, limit int, size int, totalSize int) string {

	testSingleUserRecord := `
    {
      "type": "known",
      "username": "%s",
      "userKey": "<string>",
      "accountId": "%s",
      "accountType": "atlassian",
      "email": "%s",
      "publicName": "<string>",
      "profilePicture": {
        "path": "<string>",
        "width": 2154,
        "height": 2154,
        "isDefault": true
      },
      "displayName": "<string>",
      "timeZone": "<string>",
      "isExternalCollaborator": true,
      "externalCollaborator": true,
      "operations": [
        {
          "operation": "administer",
          "targetType": "<string>"
        }
      ],
      "details": {
        "business": {
          "position": "<string>",
          "department": "<string>",
          "location": "<string>"
        },
        "personal": {
          "phone": "<string>",
          "im": "<string>",
          "website": "<string>",
          "email": "<string>"
        }
      },
      "personalSpace": {
        "id": 2154,
        "key": "<string>",
        "name": "<string>",
        "icon": {
          "path": "<string>",
          "width": 2154,
          "height": 2154,
          "isDefault": true
        },
        "description": {
          "plain": {
            "value": "<string>",
            "representation": "plain",
            "embeddedContent": [
              {}
            ]
          },
          "view": {
            "value": "<string>",
            "representation": "plain",
            "embeddedContent": [
              {}
            ]
          },
          "_expandable": {
            "view": "<string>",
            "plain": "<string>"
          }
        },
        "homepage": {
          "type": "<string>",
          "status": "<string>"
        },
        "type": "<string>",
        "metadata": {
          "labels": {
            "results": [
              {
                "prefix": "<string>",
                "name": "<string>",
                "id": "<string>",
                "label": "<string>"
              }
            ],
            "size": 2154
          },
          "_expandable": {}
        },
        "operations": [
          {
            "operation": "administer",
            "targetType": "<string>"
          }
        ],
        "permissions": [
          {
            "operation": {
              "operation": "administer",
              "targetType": "<string>"
            },
            "anonymousAccess": true,
            "unlicensedAccess": true
          }
        ],
        "status": "<string>",
        "settings": {
          "routeOverrideEnabled": true,
          "_links": {}
        },
        "theme": {
          "themeKey": "<string>"
        },
        "lookAndFeel": {
          "headings": {
            "color": "<string>"
          },
          "links": {
            "color": "<string>"
          },
          "menus": {
            "hoverOrFocus": {
              "backgroundColor": "<string>"
            },
            "color": "<string>"
          },
          "header": {
            "backgroundColor": "<string>",
            "button": {
              "backgroundColor": "<string>",
              "color": "<string>"
            },
            "primaryNavigation": {
              "color": "<string>",
              "hoverOrFocus": {
                "backgroundColor": "<string>",
                "color": "<string>"
              }
            },
            "secondaryNavigation": {
              "color": "<string>",
              "hoverOrFocus": {
                "backgroundColor": "<string>",
                "color": "<string>"
              }
            },
            "search": {
              "backgroundColor": "<string>",
              "color": "<string>"
            }
          },
          "content": {},
          "bordersAndDividers": {
            "color": "<string>"
          }
        },
        "history": {
          "createdDate": "<string>"
        },
        "_expandable": {
          "settings": "<string>",
          "metadata": "<string>",
          "operations": "<string>",
          "lookAndFeel": "<string>",
          "permissions": "<string>",
          "icon": "<string>",
          "description": "<string>",
          "theme": "<string>",
          "history": "<string>",
          "homepage": "<string>",
          "identifiers": "<string>"
        },
        "_links": {}
      },
      "_expandable": {
        "operations": "<string>",
        "details": "<string>",
        "personalSpace": "<string>"
      },
      "_links": {}
    }
`
	users := ""
	if size <= 0 {
		users = ""
	} else if size == 1 {
		id := uuid.New().String()
		users = fmt.Sprintf(testSingleUserRecord, id, id, id)
	} else {
		id := uuid.New().String()
		users = fmt.Sprintf(testSingleUserRecord, id, id, id)
		for i := 1; i < size; i++ {
			id = uuid.New().String()
			users = users + "," + fmt.Sprintf(testSingleUserRecord, id, id, id)
		}
	}

	testBody := `
    {
  "results": [%s],
  "start": %d,
  "limit": %d,
  "size": %d,
  "totalSize": %d,
  "_links": {}
}
  `
	return fmt.Sprintf(testBody, users, start, limit, size, totalSize)
}

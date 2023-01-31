package confluence

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"sort"
	"strconv"
	"strings"
)

func resourceSpacePermissionMapping() *schema.Resource {
	validPermissions := []string{
		"create:page", "create:blogpost", "create:comment", "create:attachment",
		"read:space",
		"delete:space", "delete:page", "delete:blogpost", "delete:comment", "delete:attachment",
		"export:space",
		"administer:space",
		"archive:page",
		"restrict_content:space",
	}
	return &schema.Resource{
		Create: resourceSpacePermissionMappingCreate,
		Read:   resourceSpacePermissionMappingRead,
		Delete: resourceSpacePermissionMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"operations": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(validPermissions, false),
				},
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSpacePermissionMappingCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	contentRequests := spacePermissionMappingFromResourceData(d)
	spaceKey := d.Get("key").(string)
	var createdIds []string
	for _, contentRequest := range contentRequests {
		contentResponse, err := client.CreateSpacePermission(spaceKey, contentRequest)
		if err != nil {
			return err
		}
		createdIds = append(createdIds, strconv.Itoa(contentResponse.Id))
	}
	sort.Strings(createdIds)
	d.SetId(strings.Join(createdIds[:], ":"))
	return resourceSpacePermissionMappingRead(d, m)
}

func resourceSpacePermissionMappingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	spaceKey := d.Get("key").(string)
	contentResponse, err := client.GetSpacePermission(spaceKey)
	if err != nil {
		d.SetId("")
		return err
	}
	return updateResourceDataFromSpacePermissionMapping(d, contentResponse, client)
}

func resourceSpacePermissionMappingDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	spaceKey := d.Get("key").(string)
	contentResponse, err := client.GetSpacePermission(spaceKey)
	if err != nil {
		return err
	}
	exportId := getExportSpacePermissionIdFromSpacePermissionMapping(d, contentResponse)
	for _, permissionID := range strings.Split(d.Id(), ":") {
		id, _ := strconv.Atoi(permissionID)
		err := client.DeleteSpacePermission(spaceKey, id)
		if err != nil && exportId != permissionID {
			return err
		}
	}

	// d.SetId("") is automatically called assuming delete returns no errors
	return nil
}

func spacePermissionMappingFromResourceData(d *schema.ResourceData) []*SpacePermission {
	collection := []*SpacePermission{}
	permissionsRaw := d.Get("operations").([]interface{})
	permissions := make([]string, len(permissionsRaw))
	for i, raw := range permissionsRaw {
		permissions[i] = raw.(string)
	}
	if contains(permissions, "read:space") && len(permissions) > 1 && permissions[0] != "read:space" {
		permissions = moveToFirstPositionOfSlice(permissions, "read:space")
	}
	for _, permission := range permissions {
		permissionParts := strings.Split(permission, ":")
		subject := &Subject{
			Type:       "group",
			Identifier: d.Get("group").(string),
		}
		operation := &Operation{
			Key:    permissionParts[0],
			Target: permissionParts[1],
		}
		spacePermission := &SpacePermission{
			Id:        0,
			Subject:   subject,
			Operation: operation,
		}
		collection = append(collection, spacePermission)
	}
	return collection
}

func updateResourceDataFromSpacePermissionMapping(d *schema.ResourceData, spacePermissions *SummarySpacePermissions, client *Client) error {
	var permissionIds []string
	for _, permission := range spacePermissions.Permissions {
		if permission.Subjects.Group != nil && len(permission.Subjects.Group.Results) > 0 {
			for _, group := range permission.Subjects.Group.Results {
				if d.Get("group").(string) == group.Name {
					permissionIds = append(permissionIds, strconv.Itoa(permission.ID))
				}
			}
		}
	}
	sort.Strings(permissionIds)
	d.SetId(strings.Join(permissionIds[:], ":"))
	return nil
}

func getExportSpacePermissionIdFromSpacePermissionMapping(d *schema.ResourceData, spacePermissions *SummarySpacePermissions) string {
	permissionId := ""
	for _, permission := range spacePermissions.Permissions {
		if permission.Subjects.Group != nil && len(permission.Subjects.Group.Results) > 0 {
			for _, group := range permission.Subjects.Group.Results {
				if d.Get("group").(string) == group.Name && permission.Operation.Operation == "export" && permission.Operation.TargetType == "space" {
					permissionId = strconv.Itoa(permission.ID)
				}
			}
		}
	}
	return permissionId
}

package confluence

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	contentRequest := groupFromResourceData(d)
	contentResponse, err := client.CreateGroup(contentRequest)

	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(contentResponse.Id))
	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	contentResponse, err := client.GetGroup(d.Get("id").(string))
	if err != nil {
		d.SetId("")
		return err
	}
	return updateResourceDataFromGroup(d, contentResponse, client)
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	contentRequest := groupFromResourceData(d)
	_, err := client.UpdateGroup(contentRequest)
	if err != nil {
		d.SetId("")
		return err
	}
	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := client.DeleteGroup(d.Get("id").(string))
	if err != nil {
		return err
	}
	// d.SetId("") is automatically called assuming delete returns no errors
	return nil
}

func groupFromResourceData(d *schema.ResourceData) *Group {
	id, _ := strconv.Atoi(d.Id())
	result := &Group{
		Id:   id,
		Type: "group",
		Name: d.Get("name").(string),
	}
	return result
}

func updateResourceDataFromGroup(d *schema.ResourceData, space *Group, client *Client) error {
	d.SetId(strconv.Itoa(space.Id))
	m := map[string]interface{}{
		"type": space.Type,
		"name": space.Name,
	}
	for k, v := range m {
		err := d.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

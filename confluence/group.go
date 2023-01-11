package confluence

import (
	"fmt"
	"strconv"
	"strings"
)

// Content is a primary resource in Confluence
type Group struct {
	Id    int           `json:"id,omitempty"`
	Name  string        `json:"name,omitempty"`
	Type  string        `json:"type,omitempty",default:"group"`
	Links *GenericLinks `json:"_links,omitempty"`
}

// GenericLinks is part of Content
type GenericLinks struct {
	Base  string `json:"base,omitempty"`
	WebUI string `json:"webui,omitempty"`
}

func (c *Client) CreateGroup(group *Group) (*Group, error) {
	var response Group
	if err := c.Post("/rest/api/group", group, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) GetGroup(id string) (*Group, error) {
	var response Group
	path := fmt.Sprintf("/rest/api/group/by-id?id=%s", id)
	if err := c.Get(path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) UpdateGroup(group *Group) (*Group, error) {
	if err := c.DeleteGroup(strconv.Itoa(group.Id)); err != nil {
		return nil, err
	}
	return c.CreateGroup(group)
}

func (c *Client) DeleteGroup(id string) error {
	path := fmt.Sprintf("/rest/api/group/by-id?id=%s", id)
	if err := c.Delete(path); err != nil {
		if strings.HasPrefix(err.Error(), "202 ") {
			//202 is the delete API success response
			//Other APIs return 204. Because, reasons.
			return nil
		}
		return err
	}
	return nil
}

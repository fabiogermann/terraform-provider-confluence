package confluence

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Content is a primary resource in Confluence
type SpacePermission struct {
	Id        int                   `json:"id,omitempty"`
	Subject   *Subject              `json:"subject,omitempty"`
	Operation *Operation            `json:"operation,omitempty"`
	Links     *SpacePermissionLinks `json:"_links,omitempty"`
}

// Subject is part of SpacePermission
type Subject struct {
	Type       string `json:"type,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

// Operation is part of SpacePermission
type Operation struct {
	Key    string `json:"key,omitempty"`
	Target string `json:"target,omitempty"`
}

// SpacePermissionLinks is part of SpacePermission
type SpacePermissionLinks struct {
	Base       string `json:"base,omitempty"`
	WebUI      string `json:"webui,omitempty"`
	Context    string `json:"context,omitempty"`
	Self       string `json:"self,omitempty"`
	Collection string `json:"collection,omitempty"`
}

// SummarySpacePermissions is the permissions object returned by the GET api call
type SummarySpacePermissions struct {
	ID          int               `json:"id,omitempty"`
	Key         string            `json:"key,omitempty"`
	Name        string            `json:"name,omitempty"`
	Type        string            `json:"type,omitempty"`
	Status      string            `json:"status,omitempty"`
	Permissions []SavedPermission `json:"permissions,omitempty"`

	Expandable *Expandable           `json:"_expandable,omitempty"`
	Links      *SpacePermissionLinks `json:"_links,omitempty"`
}

// SavedPermission is part of SummarySpacePermissions
type SavedPermission struct {
	ID       int `json:"id,omitempty"`
	Subjects struct {
		User struct {
			Results []struct {
				Type           string `json:"type,omitempty"`
				AccountID      string `json:"accountId,omitempty"`
				AccountType    string `json:"accountType,omitempty"`
				Email          string `json:"email,omitempty"`
				PublicName     string `json:"publicName,omitempty"`
				ProfilePicture struct {
					Path      string `json:"path,omitempty"`
					Width     int    `json:"width,omitempty"`
					Height    int    `json:"height,omitempty"`
					IsDefault bool   `json:"isDefault,omitempty"`
				} `json:"profilePicture,omitempty"`
				DisplayName            string                `json:"displayName,omitempty"`
				IsExternalCollaborator bool                  `json:"isExternalCollaborator,omitempty"`
				Expandable             *Expandable           `json:"_expandable,omitempty"`
				Links                  *SpacePermissionLinks `json:"_links,omitempty"`
			} `json:"results,omitempty"`
			Size int `json:"size,omitempty"`
		} `json:"user,omitempty"`
		Group      *SavedPermissionGroup `json:"group,omitempty"`
		Expandable *Expandable           `json:"_expandable,omitempty"`
	} `json:"subjects,omitempty"`
	Operation struct {
		Operation  string `json:"operation,omitempty"`
		TargetType string `json:"targetType,omitempty"`
	} `json:"operation,omitempty"`
	AnonymousAccess  bool `json:"anonymousAccess,omitempty"`
	UnlicensedAccess bool `json:"unlicensedAccess,omitempty"`
}

// SavedPermissionGroup is part of SavedPermission
type SavedPermissionGroup struct {
	Results []struct {
		Type  string                `json:"type,omitempty"`
		Name  string                `json:"name,omitempty"`
		ID    string                `json:"id,omitempty"`
		Links *SpacePermissionLinks `json:"_links,omitempty"`
	} `json:"results,omitempty"`
	Size int `json:"size,omitempty"`
}

// Expandable is part of SummarySpacePermissions
type Expandable struct {
	Settings      string `json:"settings,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
	Operations    string `json:"operations,omitempty"`
	PersonalSpace string `json:"personalSpace,omitempty"`
	LookAndFeel   string `json:"lookAndFeel,omitempty"`
	Identifiers   string `json:"identifiers,omitempty"`
	Icon          string `json:"icon,omitempty"`
	Description   string `json:"description,omitempty"`
	Theme         string `json:"theme,omitempty"`
	History       string `json:"history,omitempty"`
	Homepage      string `json:"homepage,omitempty"`
	User          string `json:"user,omitempty"`
}

func (c *Client) CreateSpacePermission(spaceKey string, spacePermission *SpacePermission) (*SpacePermission, error) {
	var response SpacePermission
	path := fmt.Sprintf("/rest/api/space/%s/permission", spaceKey)

	b, _ := json.Marshal(&spacePermission)
	log.Printf("[INFO] >>>>>>>>>>> %s", b)
	if err := c.Post(path, spacePermission, &response); err != nil {
		log.Printf("[INFO] >>>>>>>>>>> %s", err.Error())
		return nil, err
	}
	return &response, nil
}

func (c *Client) GetSpacePermission(spaceKey string) (*SummarySpacePermissions, error) {
	var response SummarySpacePermissions
	path := fmt.Sprintf("/rest/api/space/%s?expand=permissions", spaceKey)
	if err := c.Get(path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) UpdateSpacePermission(spaceKey string, spacePermission *SpacePermission) (*SpacePermission, error) {
	if err := c.DeleteSpacePermission(spaceKey, spacePermission.Id); err != nil {
		return nil, err
	}
	return c.CreateSpacePermission(spaceKey, spacePermission)
}

func (c *Client) DeleteSpacePermission(spaceKey string, id int) error {
	path := fmt.Sprintf("/rest/api/space/%s/permission/%d", spaceKey, id)
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

package transferobjects

// Response Object
type GroupMembersResponse struct {
	Members   []Member           `json:"results,omitempty"`
	Start     int                `json:"start,omitempty"`
	Limit     int                `json:"limit,omitempty"`
	Size      int                `json:"size,omitempty"`
	TotalSize int                `json:"totalSize,omitempty"`
	Links     *GroupMembersLinks `json:"_links,omitempty"`
}

type AccountIDRecord struct {
	AccountID string `json:"accountId,omitempty"`
}

type Member struct {
	Type           string `json:"type,omitempty"`
	Username       string `json:"username,omitempty"`
	UserKey        string `json:"userKey,omitempty"`
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
	DisplayName            string `json:"displayName,omitempty"`
	TimeZone               string `json:"timeZone,omitempty"`
	IsExternalCollaborator bool   `json:"isExternalCollaborator,omitempty"`
	ExternalCollaborator   bool   `json:"externalCollaborator,omitempty"`
	Expandable             struct {
		Operations    string `json:"operations,omitempty"`
		PersonalSpace string `json:"personalSpace,omitempty"`
	} `json:"_expandable,omitempty"`
	Links struct {
		Self string `json:"self,omitempty"`
	} `json:"_links,omitempty"`
}

// GenericLinks is part of Content
type GroupMembersLinks struct {
	Base    string `json:"base,omitempty"`
	WebUI   string `json:"webui,omitempty"`
	Context string `json:"context,omitempty"`
	Self    string `json:"self,omitempty"`
}

package transferobjects

// Content is a primary resource in Confluence
type Group struct {
	Id    string      `json:"id,omitempty"`
	Name  string      `json:"name,omitempty"`
	Type  string      `json:"type,omitempty"`
	Links *GroupLinks `json:"_links,omitempty"`
}

// GenericLinks is part of Content
type GroupLinks struct {
	Base  string `json:"base,omitempty"`
	WebUI string `json:"webui,omitempty"`
}

package transferobjects

// Content is a primary resource in Confluence
type Space struct {
	Id    FlexInt     `json:"id,omitempty"`
	Name  string      `json:"name,omitempty"`
	Key   string      `json:"key,omitempty"`
	Links *SpaceLinks `json:"_links,omitempty"`
}

// ContentLinks is part of Content
type SpaceLinks struct {
	Base  string `json:"base,omitempty"`
	WebUI string `json:"webui,omitempty"`
}

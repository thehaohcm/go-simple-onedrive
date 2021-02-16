package models

type LinkShared struct {
	Type   string `json:"type,omitempty"`
	Scope  string `json:"scope,omitempty"`
	WebUrl string `json:"webUrl,omitempty"`
}

package models

type SharedLinkResponse struct {
	Id    string     `json:"id,omitempty"`
	Roles []string   `json:"roles,omitempty"`
	Link  LinkShared `json:"link,omitempty"`
}
